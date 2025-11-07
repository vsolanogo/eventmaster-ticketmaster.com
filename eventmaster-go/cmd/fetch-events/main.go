package main

import (
	"eventmaster-go/internal/config"
	"eventmaster-go/internal/database"
	"eventmaster-go/internal/models"
	"eventmaster-go/internal/repositories"
	"eventmaster-go/internal/services"
	"flag"
	"log"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "../../.env", "path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
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
			log.Printf("Failed to close database: %v", err)
		}
	}()

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

	// Initialize repositories
	eventRepo := repositories.NewEventRepository(db)
	userRepo := repositories.NewUserRepository(db)
	participantRepo := repositories.NewParticipantRepository(db)
	imageRepo := repositories.NewImageRepository(db)

	// Prepare dependencies
	imageService := services.NewImageService(imageRepo)
	participantService := services.NewParticipantService(participantRepo, eventRepo)
	systemUserID, err := services.EnsureTicketmasterSystemUser(userRepo)
	if err != nil {
		log.Fatalf("Failed to ensure Ticketmaster system user: %v", err)
	}

	// Initialize services
	ticketmasterService := services.NewTicketmasterService(
		eventRepo,
		imageService,
		participantService,
		cfg.Ticketmaster.APIKey,
		systemUserID,
	)

	// Fetch and save events
	log.Println("Fetching events from Ticketmaster...")
	if err := ticketmasterService.FetchAndSaveEvents(); err != nil {
		log.Fatalf("Failed to fetch events: %v", err)
	}

	log.Println("Successfully fetched and saved events from Ticketmaster")
}
