package main

import (
	"context"
	"log"
	"os"
	"time"

	"eventmaster-go/internal/config"
	"eventmaster-go/internal/database"
	"eventmaster-go/internal/models"
	"eventmaster-go/internal/repositories"
	"eventmaster-go/internal/server"
	"eventmaster-go/internal/services"

	"gorm.io/gorm"
)

// Config is now defined in the config package

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("../../.env")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.NewDB(&database.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.Username,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := database.CloseDB(db); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Auto-migrate the schema
	if err := db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Event{},
		&models.Participant{},
		&models.Image{},
		&models.Session{},
	); err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	// Seed initial data
	if err := seedInitialData(db, cfg.Auth.AdminEmail, cfg.Auth.AdminPassword); err != nil {
		log.Fatalf("Failed to seed initial data: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	eventRepo := repositories.NewEventRepository(db)
	participantRepo := repositories.NewParticipantRepository(db)
	imageRepo := repositories.NewImageRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)

	// Initialize services
	authService := services.NewAuthService(
		userRepo,
		sessionRepo,
		cfg.Auth.JWTExpiration,
	)
	eventService := services.NewEventService(eventRepo, imageRepo)
	participantService := services.NewParticipantService(participantRepo, eventRepo)
	imageService := services.NewImageService(imageRepo)
	systemUserID, err := services.EnsureTicketmasterSystemUser(userRepo)
	if err != nil {
		log.Fatalf("Failed to ensure Ticketmaster system user: %v", err)
	}
	ticketmasterService := services.NewTicketmasterService(
		eventRepo,
		imageService,
		participantService,
		cfg.Ticketmaster.APIKey,
		systemUserID,
	)

	schedulerCtx, schedulerCancel := context.WithCancel(context.Background())
	defer schedulerCancel()
	ticketmasterService.StartScheduler(schedulerCtx, 6*time.Hour)

	go func() {
		const initialFetchDelay = 5 * time.Second
		log.Printf("Delaying Ticketmaster fetch until %s after server start", initialFetchDelay)
		time.Sleep(initialFetchDelay)
		if err := ticketmasterService.FetchAndSaveEvents(); err != nil {
			log.Printf("Ticketmaster fetch after delay failed: %v", err)
			return
		}
		log.Println("Ticketmaster fetch completed after delayed start")
	}()

	// Configure file upload paths
	uploadPath := "./uploads"
	baseURL := "/uploads/"
	
	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}
	
	fileService := services.NewFileService(imageRepo, uploadPath, baseURL)

	// Initialize HTTP server
	serverPort := "3000"
	serverConfig := &server.Config{
		Port:              serverPort,
		SessionCookieName: cfg.Server.SessionCookieName,
	}

	srv := server.NewServer(authService, *serverConfig)

	// Register handlers
	srv.RegisterEventHandlers(eventService)
	srv.RegisterParticipantHandlers(participantService)
	srv.RegisterFileHandlers(fileService)

	log.Printf("Server starting on :%s\n", serverPort)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// seedInitialData populates the database with initial required data
func seedInitialData(db *gorm.DB, adminEmail, adminPassword string) error {
	// Check if admin role exists
	var adminRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create admin role
			adminRole = models.Role{
				Name: "admin",
			}
			if err := db.Create(&adminRole).Error; err != nil {
				log.Printf("Failed to create admin role: %v", err)
				return err
			}

			// Create user role
			userRole := models.Role{
				Name: "user",
			}
			if err := db.Create(&userRole).Error; err != nil {
				log.Printf("Failed to create user role: %v", err)
				return err
			}
		} else {
			log.Printf("Failed to check for admin role: %v", err)
			return err
		}
	}

	// Check if admin user exists
	var adminUser models.User
	if err := db.Where("email = ?", adminEmail).First(&adminUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create admin user
			adminUser = models.User{
				Email:    adminEmail,
				Password: adminPassword,
			}
			if err := db.Create(&adminUser).Error; err != nil {
				log.Printf("Failed to create admin user: %v", err)
				return err
			}

			// Assign admin role to admin user
			if err := db.Model(&adminUser).Association("Roles").Append(&adminRole); err != nil {
				log.Printf("Failed to assign admin role: %v", err)
				return err
			}
		} else {
			log.Printf("Failed to check for admin user: %v", err)
			return err
		}
	}

	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
