package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"eventmaster-go/internal/models"
	"eventmaster-go/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

const defaultUserRole = "user"

// AuthService handles authentication-related business logic
type AuthService interface {
	Register(email, password string) (*models.User, error)
	Login(email, password, ip string) (*models.Session, error)
	GetUserByID(id string) (*models.User, error)
	ValidateSession(token string) (*models.User, error)
	Logout(token string) error
}

type authService struct {
	userRepo    repositories.UserRepository
	sessionRepo repositories.SessionRepository
	tokenExpiry time.Duration
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repositories.UserRepository, sessionRepo repositories.SessionRepository, tokenExpiry time.Duration) AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		tokenExpiry: tokenExpiry,
	}
}

func (s *authService) Register(email, password string) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create new user. The User model's BeforeCreate hook will hash the password.
	user := &models.User{
		Email:    email,
		Password: password,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	if err := s.userRepo.AttachRoleByName(user, defaultUserRole); err != nil {
		return nil, err
	}

	return s.userRepo.FindWithAssociations(user.ID)
}

func (s *authService) Login(email, password, ip string) (*models.Session, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := generateSessionToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(s.tokenExpiry)

	session := &models.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
		IP:        ip,
		User:      *user,
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *authService) GetUserByID(id string) (*models.User, error) {
	return s.userRepo.FindWithAssociations(id)
}

func (s *authService) ValidateSession(token string) (*models.User, error) {
	session, err := s.sessionRepo.FindByToken(token)
	if err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		_ = s.sessionRepo.DeleteByToken(token)
		return nil, errors.New("session expired")
	}

	return &session.User, nil
}

func (s *authService) Logout(token string) error {
	if token == "" {
		return nil
	}
	return s.sessionRepo.DeleteByToken(token)
}

func generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
