package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"eventmaster-go/internal/config"
	"eventmaster-go/internal/database"
	"eventmaster-go/internal/models"
	"eventmaster-go/internal/repositories"
	"eventmaster-go/internal/server"
	"eventmaster-go/internal/services"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	testClient        *http.Client
	apiBaseURL        string
	sessionCookieName string
	srv               *server.Server
	db                *gorm.DB
	cancelFn          context.CancelFunc
	subtestsPassed    atomic.Int32
	subtestsFailed    atomic.Int32
	resultsMu         sync.Mutex
	subtestResults    []testResult
)

type testResult struct {
	Name   string `json:"name"`
	Passed bool   `json:"passed"`
}

type summaryReport struct {
	Total    int          `json:"total"`
	Passed   int          `json:"passed"`
	Failed   int          `json:"failed"`
	Subtests []testResult `json:"subtests"`
}

func runSubtest(t *testing.T, name string, fn func(t *testing.T)) {
	fullName := fmt.Sprintf("%s/%s", t.Name(), name)
	if t.Run(name, fn) {
		subtestsPassed.Add(1)
		recordSubtestResult(fullName, true)
	} else {
		subtestsFailed.Add(1)
		recordSubtestResult(fullName, false)
	}
}

func recordSubtestResult(name string, passed bool) {
	resultsMu.Lock()
	defer resultsMu.Unlock()
	subtestResults = append(subtestResults, testResult{Name: name, Passed: passed})
}

func writeSummaryReport() error {
	resultsMu.Lock()
	resultsCopy := make([]testResult, len(subtestResults))
	copy(resultsCopy, subtestResults)
	resultsMu.Unlock()

	report := summaryReport{
		Total:    len(resultsCopy),
		Passed:   int(subtestsPassed.Load()),
		Failed:   int(subtestsFailed.Load()),
		Subtests: resultsCopy,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	resultsPath, err := resolveResultsFilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(resultsPath), 0o755); err != nil {
		return err
	}

	if err := os.WriteFile(resultsPath, data, 0o644); err != nil {
		return err
	}

	return nil
}

func resolveResultsFilePath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to determine caller for results path")
	}
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "test_results.json"), nil
}

func TestMain(m *testing.M) {
	// Load configuration
	cfg, err := config.LoadConfig("../.env")
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// Force the server to listen on a dedicated port for tests
	cfg.Server.Port = getEnv("E2E_TEST_PORT", "3100")
	serverURL := fmt.Sprintf("http://localhost:%s", cfg.Server.Port)
	apiBaseURL = serverURL + "/api"
	sessionCookieName = cfg.Server.SessionCookieName

	// Initialise database connection
	db, err = database.NewDB(&database.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.Username,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect DB: %v", err))
	}

	// Auto-migrate the schema to ensure tables exist
	if err := db.AutoMigrate(&models.Role{}, &models.User{}, &models.Event{}, &models.Participant{}, &models.Image{}, &models.Session{}); err != nil {
		panic(fmt.Sprintf("failed to automigrate: %v", err))
	}

	// Set up repositories
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	eventRepo := repositories.NewEventRepository(db)
	participantRepo := repositories.NewParticipantRepository(db)
	imageRepo := repositories.NewImageRepository(db)

	// Set up services
	authService := services.NewAuthService(userRepo, sessionRepo, cfg.Auth.JWTExpiration)
	eventService := services.NewEventService(eventRepo, imageRepo)
	participantService := services.NewParticipantService(participantRepo, eventRepo)
	imageService := services.NewImageService(imageRepo)
	fileService := services.NewFileService(imageRepo, "./uploads", "/uploads")
	systemUserID, err := services.EnsureTicketmasterSystemUser(userRepo)
	if err != nil {
		panic(fmt.Sprintf("failed to ensure Ticketmaster system user: %v", err))
	}

	ticketmasterService := services.NewTicketmasterService(
		eventRepo,
		imageService,
		participantService,
		cfg.Ticketmaster.APIKey,
		systemUserID,
	)

	serverConfig := server.Config{
		Port:              cfg.Server.Port,
		SessionCookieName: cfg.Server.SessionCookieName,
	}

	srv = server.NewServer(authService, serverConfig)
	srv.RegisterEventHandlers(eventService)
	srv.RegisterParticipantHandlers(participantService)
	srv.RegisterFileHandlers(fileService)

	srvCtx, cancel := context.WithCancel(context.Background())
	cancelFn = cancel

	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Sprintf("server failed to start: %v", err))
		}
	}()

	go ticketmasterService.StartScheduler(srvCtx, time.Hour)

	waitForServer(serverURL)

	jar, _ := cookiejar.New(nil)
	testClient = &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}

	code := m.Run()
	if err := writeSummaryReport(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write e2e summary file: %v\n", err)
	}
	fmt.Fprintf(os.Stderr, "E2E summary: passed=%d failed=%d\n", subtestsPassed.Load(), subtestsFailed.Load())

	cancelFn()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
	_ = database.CloseDB(db)

	os.Exit(code)
}

func waitForServer(base string) {
	retries := 0
	for {
		if retries > 30 {
			panic("server did not become ready in time")
		}
		time.Sleep(200 * time.Millisecond)
		resp, err := http.Get(base + "/api/health")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		retries++
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

// Helper and data structures

type UserResponse struct {
	ID      string         `json:"id"`
	Email   string         `json:"email"`
	Roles   []RoleResponse `json:"role"`
	Session []SessionInfo  `json:"session"`
	Created time.Time      `json:"createdAt"`
	Updated time.Time      `json:"updatedAt"`
}

type RoleResponse struct {
	Role        string  `json:"role"`
	Description *string `json:"description"`
}

type SessionInfo struct {
	ID        string    `json:"id"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"createdAt"`
	Expires   time.Time `json:"expires"`
}

type EventListResponse struct {
	Events     []EventResponse `json:"events"`
	TotalCount int             `json:"totalCount"`
}

type EventResponse struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	EventDate *time.Time `json:"eventDate"`
}

type ParticipantResponse struct {
	ID                string     `json:"id"`
	FullName          string     `json:"fullName"`
	Email             string     `json:"email"`
	DateOfBirth       *time.Time `json:"dateOfBirth"`
	SourceOfDiscovery string     `json:"sourceOfDiscovery"`
	EventID           string     `json:"eventId"`
}

func createUser(t *testing.T, email, password string) string {
	resp := doRequest(t, http.MethodPost, "/register", map[string]string{
		"email":    email,
		"password": password,
	}, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("failed to create user: status=%d body=%s", resp.StatusCode, string(body))
	}

	var user UserResponse
	decodeJSON(t, resp.Body, &user)
	return user.ID
}

func loginAndGetCookie(t *testing.T, email, password string) string {
	createUser(t, email, password)
	resp := doRequest(t, http.MethodPost, "/login", map[string]string{
		"email":    email,
		"password": password,
	}, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected login success, got %d", resp.StatusCode)
	}

	cookies := resp.Cookies()
	if len(cookies) == 0 {
		t.Fatalf("expected session cookie")
	}

	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if strings.EqualFold(c.Name, sessionCookieName) {
			sessionCookie = c
			break
		}
	}

	if sessionCookie == nil {
		var names []string
		for _, c := range cookies {
			names = append(names, c.Name)
		}
		t.Fatalf("expected cookie %s, got %v", sessionCookieName, names)
	}

	return fmt.Sprintf("%s=%s", sessionCookie.Name, sessionCookie.Value)
}

func doRequest(t *testing.T, method, path string, body any, headers map[string]string) *http.Response {
	t.Helper()

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal body: %v", err)
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, apiBaseURL+path, reader)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := testClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	return resp
}

func decodeJSON(t *testing.T, r io.Reader, v any) {
	t.Helper()
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(v); err != nil && !errors.Is(err, io.EOF) {
		t.Fatalf("failed to decode JSON: %v", err)
	}
}

func resetCookies(t *testing.T) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("failed to reset cookies: %v", err)
	}
	testClient.Jar = jar
}

func randomEmail() string {
	uid := uuid.NewString()
	return fmt.Sprintf("test-%s@example.com", uid[:8])
}

func requireEventID(t *testing.T) string {
	events := fetchEvents(t)
	if len(events) == 0 {
		t.Skip("no events available")
	}
	return events[0].ID
}

func fetchEvents(t *testing.T) []EventResponse {
	resp := doRequest(t, http.MethodGet, "/events", nil, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var list EventListResponse
	decodeJSON(t, resp.Body, &list)
	return list.Events
}

func registerParticipant(t *testing.T, eventID, fullName, email string) ParticipantResponse {
	payload := map[string]any{
		"fullName":          fullName,
		"email":             email,
		"dateOfBirth":       time.Now().AddDate(-30, 0, 0).Format(time.RFC3339),
		"sourceOfDiscovery": "social_media",
	}

	resp := doRequest(t, http.MethodPost, fmt.Sprintf("/participants?eventId=%s", eventID), payload, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status %d, got %d body=%s", http.StatusCreated, resp.StatusCode, string(body))
	}

	var participant ParticipantResponse
	decodeJSON(t, resp.Body, &participant)
	return participant
}
