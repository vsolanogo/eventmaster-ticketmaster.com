package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	DB         DBConfig
	Server     ServerConfig
	Auth       AuthConfig
	Ticketmaster TicketmasterConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	SSLMode  string
}

type ServerConfig struct {
	Port             string
	SessionCookieName string
}

type AuthConfig struct {
	JWTSecret     string
	JWTExpiration time.Duration
	SessionID     string
	AdminEmail    string
	AdminPassword string
}

type TicketmasterConfig struct {
	APIKey string
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig(envPath string) (*Config, error) {
	// First try to load from the current directory
	if err := godotenv.Load(); err != nil {
		// If that fails, try to load from the provided path
		if err := godotenv.Load(envPath); err != nil {
			log.Printf("Warning: Could not load .env file: %v", err)
		}
	}

	// Explicitly load the Ticketmaster key from environment
	tmKey := os.Getenv("TICKETMASTER_KEY")

	config := &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5433"),
			Username: getEnv("DB_USERNAME", "postgres"),
			Password: getEnv("DB_PASSWORD", "passwordSuperUser1111"),
			Name:     getEnv("DB_NAME", "events"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port:             getEnv("SERVER_PORT", "3000"),
			SessionCookieName: getEnv("SESSION_ID", "SessionID"),
		},
		Auth: AuthConfig{
			JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-here"),
			JWTExpiration: 24 * time.Hour,
			SessionID:     getEnv("SESSION_ID", "SessionID"),
			AdminEmail:    getEnv("ROOT_ADMIN_EMAIL", "admin@admin.com"),
			AdminPassword: getEnv("ROOT_ADMIN_PASSWORD", "admin"),
		},
		Ticketmaster: TicketmasterConfig{
			APIKey: tmKey,
		},
	}

	// Validate required configurations
	if config.Ticketmaster.APIKey == "" {
		log.Println("WARNING: TICKETMASTER_KEY is not set. Ticketmaster integration will not work.")
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
