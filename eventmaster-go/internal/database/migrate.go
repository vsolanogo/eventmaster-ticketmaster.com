package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/gorm"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// GetMigrationFiles returns a sorted list of migration files
func GetMigrationFiles() ([]string, error) {
	files, err := os.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrations []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			migrations = append(migrations, filepath.Join("migrations", file.Name()))
		}
	}

	return migrations, nil
}

// RunMigrations runs database migrations
func RunMigrations(db *gorm.DB) error {
	// Get the underlying sql.DB instance
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Configure goose
	goose.SetTableName("goose_db_version")

	// Get the absolute path to the migrations directory
	migrationsDir := filepath.Join("migrations")

	// Ensure the migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", migrationsDir)
	}

	// Run migrations
	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Verify the migration
	currentVersion, err := goose.GetDBVersion(sqlDB)
	if err != nil {
		return fmt.Errorf("failed to get current migration version: %w", err)
	}

	log.Printf("Database migrations completed successfully. Current version: %d\n", currentVersion)
	return nil
}

// RollbackMigrations rolls back the database schema by one version
func RollbackMigrations(databaseURL string) error {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Get current version
	currentVersion, err := goose.GetDBVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current migration version: %w", err)
	}

	if currentVersion == 0 {
		log.Println("No migrations to roll back")
		return nil
	}

	// Run rollback
	if err := goose.Down(db, "./migrations"); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Printf("Successfully rolled back to version %d\n", currentVersion-1)
	return nil
}

