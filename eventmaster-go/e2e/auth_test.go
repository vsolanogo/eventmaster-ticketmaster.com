package e2e

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestAuthFlow(t *testing.T) {
	runSubtest(t, "register success", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    randomEmail(),
			"password": "StrongPassw0rd!",
		}

		resp := doRequest(t, http.MethodPost, "/register", reqBody, nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}

		var userResp UserResponse
		decodeJSON(t, resp.Body, &userResp)

		if userResp.ID == "" {
			t.Fatalf("expected user ID to be set")
		}
	})

	runSubtest(t, "register duplicate", func(t *testing.T) {
		email := randomEmail()
		password := "Duplicate1!"
		createUser(t, email, password)

		resp := doRequest(t, http.MethodPost, "/register", map[string]string{
			"email":    email,
			"password": password,
		}, nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})

	runSubtest(t, "login success", func(t *testing.T) {
		email := randomEmail()
		password := "ValidPass1!"
		createUser(t, email, password)

		resp := doRequest(t, http.MethodPost, "/login", map[string]string{
			"email":    email,
			"password": password,
		}, nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		if len(resp.Cookies()) == 0 {
			t.Fatalf("expected session cookie to be set")
		}
	})

	runSubtest(t, "login invalid credentials", func(t *testing.T) {
		resp := doRequest(t, http.MethodPost, "/login", map[string]string{
			"email":    "unknown@example.com",
			"password": "wrong",
		}, nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
		}
	})

	runSubtest(t, "current user requires auth", func(t *testing.T) {
		email := randomEmail()
		password := "ValidPass2!"
		cookie := loginAndGetCookie(t, email, password)

		headers := map[string]string{"Cookie": cookie}
		resp := doRequest(t, http.MethodGet, "/user", nil, headers)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var user UserResponse
		decodeJSON(t, resp.Body, &user)

		if user.ID == "" {
			t.Fatalf("expected user id in response")
		}
		if !strings.EqualFold(user.Email, email) {
			t.Fatalf("expected email %s, got %s", email, user.Email)
		}
		if len(user.Roles) == 0 {
			t.Fatalf("expected at least one role, got %d", len(user.Roles))
		}
		for _, r := range user.Roles {
			if r.Role == "" {
				t.Fatalf("role missing role field: %+v", r)
			}
			if r.Description != nil && *r.Description == "" {
				t.Fatalf("role missing description field: %+v", r)
			}
		}
		if len(user.Session) == 0 {
			t.Fatalf("expected at least one active session, got %d", len(user.Session))
		}
		for _, s := range user.Session {
			if s.ID == "" {
				t.Fatalf("session missing id: %+v", s)
			}
			if s.IP == "" {
				t.Fatalf("session missing ip: %+v", s)
			}
		}
	})

	runSubtest(t, "current user unauthorized", func(t *testing.T) {
		resetCookies(t)
		resp := doRequest(t, http.MethodGet, "/user", nil, nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
		}
	})

	runSubtest(t, "current user invalid token", func(t *testing.T) {
		headers := map[string]string{
			"Cookie": fmt.Sprintf("%s=%s", sessionCookieName, "invalid-token"),
		}
		resp := doRequest(t, http.MethodGet, "/user", nil, headers)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
		}
	})
}

func TestAuthValidations(t *testing.T) {
	runSubtest(t, "register invalid email", func(t *testing.T) {
		payload := map[string]any{
			"email":    "not-an-email",
			"password": "ValidPassword123",
		}
		resp := doRequest(t, http.MethodPost, "/register", payload, nil)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})

	runSubtest(t, "register missing password", func(t *testing.T) {
		payload := map[string]any{
			"email": randomEmail(),
		}
		resp := doRequest(t, http.MethodPost, "/register", payload, nil)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})

	runSubtest(t, "register empty email", func(t *testing.T) {
		payload := map[string]any{
			"email":    "",
			"password": "ValidPassword123",
		}
		resp := doRequest(t, http.MethodPost, "/register", payload, nil)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})
}

func TestLoginEdgeCases(t *testing.T) {
	runSubtest(t, "login wrong password", func(t *testing.T) {
		email := randomEmail()
		correctPassword := "CorrectPass1!"
		createUser(t, email, correctPassword)

		resp := doRequest(t, http.MethodPost, "/login", map[string]any{
			"email":    email,
			"password": "WrongPassword123",
		}, nil)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
		}
	})

	runSubtest(t, "login missing credentials", func(t *testing.T) {
		resp := doRequest(t, http.MethodPost, "/login", map[string]any{}, nil)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})
}

func TestLogout(t *testing.T) {
	runSubtest(t, "logout clears session cookie", func(t *testing.T) {
		email := randomEmail()
		password := "LogoutPass1!"
		cookie := loginAndGetCookie(t, email, password)

		headers := map[string]string{"Cookie": cookie}
		resp := doRequest(t, http.MethodPost, "/logout", nil, headers)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}

		cookies := resp.Cookies()
		var cleared bool
		for _, c := range cookies {
			if strings.EqualFold(c.Name, sessionCookieName) {
				if c.Value == "" {
					cleared = true
					break
				}
			}
		}

		if !cleared {
			t.Fatalf("expected %s cookie to be cleared", sessionCookieName)
		}
	})
}
