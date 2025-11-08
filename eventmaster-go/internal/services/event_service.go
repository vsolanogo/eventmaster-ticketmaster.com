package services

import (
	"eventmaster-go/internal/models"
	"eventmaster-go/internal/repositories"
	"time"
	"errors"
)

// EventService handles event-related business logic
type EventService interface {
	CreateEvent(event *models.Event, userID string, imageIDs []string) (*models.Event, error)
	GetEventByID(id string) (*models.Event, error)
	GetEventsByDateRange(start, end time.Time) ([]*models.Event, error)
	GetUserEvents(userID string) ([]*models.Event, error)
	GetPaginatedEvents(page, limit int, sortBy, sortOrder string) ([]*models.Event, int64, error)
	UpdateEvent(id string, event *models.Event) (*models.Event, error)
	DeleteEvent(id string) error
}

type eventService struct {
	eventRepo repositories.EventRepository
	imageRepo repositories.ImageRepository
}

// NewEventService creates a new event service
func NewEventService(eventRepo repositories.EventRepository, imageRepo repositories.ImageRepository) EventService {
	return &eventService{
		eventRepo: eventRepo,
		imageRepo: imageRepo,
	}
}

var ErrImageNotFound = errors.New("one or more images were not found")

func (s *eventService) CreateEvent(event *models.Event, userID string, imageIDs []string) (*models.Event, error) {
	event.UserID = userID
	
	var images []*models.Image
	var err error
	if len(imageIDs) > 0 {
		images, err = s.imageRepo.FindByIDs(imageIDs)
		if err != nil {
			return nil, err
		}
		if len(images) != len(imageIDs) {
			return nil, ErrImageNotFound
		}
	}

	if err := s.eventRepo.Create(event); err != nil {
		return nil, err
	}

	if len(images) > 0 {
		if err := s.imageRepo.AttachToEvent(event, images); err != nil {
			return nil, err
		}
	}

	// Return the created event with related data
	return s.eventRepo.FindWithImages(event.ID)
}

func (s *eventService) GetEventByID(id string) (*models.Event, error) {
	return s.eventRepo.FindWithImages(id)
}

func (s *eventService) GetEventsByDateRange(start, end time.Time) ([]*models.Event, error) {
	return s.eventRepo.FindByDateRange(start, end)
}

func (s *eventService) GetUserEvents(userID string) ([]*models.Event, error) {
	return s.eventRepo.FindByUserID(userID)
}

func (s *eventService) GetPaginatedEvents(page, limit int, sortBy, sortOrder string) ([]*models.Event, int64, error) {
	return s.eventRepo.FindPaginated(page, limit, sortBy, sortOrder)
}

func (s *eventService) UpdateEvent(id string, event *models.Event) (*models.Event, error) {
	existingEvent, err := s.eventRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	existingEvent.Title = event.Title
	existingEvent.Description = event.Description
	existingEvent.Organizer = event.Organizer
	existingEvent.EventDate = event.EventDate
	existingEvent.Latitude = event.Latitude
	existingEvent.Longitude = event.Longitude

	if err := s.eventRepo.Update(existingEvent); err != nil {
		return nil, err
	}

	return s.eventRepo.FindWithImages(id)
}

func (s *eventService) DeleteEvent(id string) error {
	return s.eventRepo.Delete(id)
}
