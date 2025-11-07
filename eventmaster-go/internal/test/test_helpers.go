package test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"

	"github.com/google/uuid"
)

type TestUser struct {
	ID       string
	Email    string
	Password string
	Token    string
}

// GenerateRandomEmail generates a random email for testing
func GenerateRandomEmail() string {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	return fmt.Sprintf("test-%s@example.com", hex.EncodeToString(randomBytes))
}

// GenerateRandomPassword generates a random password for testing
func GenerateRandomPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, 16)
	for i := range password {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[n.Int64()]
	}
	return string(password)
}

// ExtractTokenFromResponse extracts the JWT token from the Set-Cookie header
func ExtractTokenFromResponse(resp *http.Response) string {
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
			return cookie.Value
		}
	}
	return ""
}

// GenerateRandomString generates a random string of the given length
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

// GenerateUUID generates a new UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// CleanupTestUser deletes a test user by email
func CleanupTestUser(email string, client *http.Client, baseURL string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/users/%s", baseURL, email), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("failed to delete test user: %s", resp.Status)
	}

	return nil
}
