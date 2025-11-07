package database

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds database configuration
type Config struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxIdleConns int
	MaxOpenConns int
}

// NewConfigFromEnv creates a new Config from environment variables
func NewConfigFromEnv() (*Config, error) {
	cfg := &Config{
		Host:         getEnv("DB_HOST", "localhost"),
		Port:         getEnv("DB_PORT", "5433"),
		User:         getEnv("DB_USER", "postgres"),
		Password:     getEnv("DB_PASSWORD", "passwordSuperUser1111"),
		DBName:       getEnv("DB_NAME", "events"),
		SSLMode:      getEnv("DB_SSLMODE", "disable"),
		MaxIdleConns: 10,
		MaxOpenConns: 100,
	}

	// Required fields validation
	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("missing required database configuration")
	}

	return cfg, nil
}

// ConnectionString returns the connection string for PostgreSQL
func (c *Config) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
		c.SSLMode,
	)
}

// NewDB creates a new database connection and sets up the schema
func NewDB(cfg *Config) (*gorm.DB, error) {
	// First connect without the database name to drop and recreate it
	dropDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.SSLMode,
	)

	// Connect to postgres database to drop and recreate our database
	db, err := gorm.Open(postgres.Open(dropDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Drop the database if it exists
	dropDB := fmt.Sprintf("DROP DATABASE IF EXISTS %s", cfg.DBName)
	if err := db.Exec(dropDB).Error; err != nil {
		return nil, fmt.Errorf("failed to drop database: %w", err)
	}

	// Create the database
	createDB := fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)
	if err := db.Exec(createDB).Error; err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	// Now connect to the new database
	dsn := cfg.ConnectionString()
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable uuid-ossp extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return nil, fmt.Errorf("failed to enable uuid-ossp extension: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// CloseDB closes the database connection
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %w", err)
	}
	return sqlDB.Close()
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
