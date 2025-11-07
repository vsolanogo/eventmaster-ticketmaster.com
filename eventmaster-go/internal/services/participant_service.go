package services

import (
	"eventmaster-go/internal/models"
	"eventmaster-go/internal/repositories"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// ParticipantService handles participant-related business logic
type ParticipantService interface {
	RegisterParticipant(participant *models.Participant) (*models.Participant, error)
	GetEventParticipants(eventID string) ([]*models.Participant, error)
	GetParticipantByID(id string) (*models.Participant, error)
	GetParticipantByEmail(email string) ([]*models.Participant, error)
	GetEventParticipantCount(eventID string) (int64, error)
	DeleteParticipant(id string) error
	GenerateFakeParticipants(event *models.Event, count int) error
}

func (s *participantService) GenerateFakeParticipants(event *models.Event, count int) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}

	if count <= 0 {
		return nil
	}

	participants := make([]models.Participant, 0, count)
	for i := 0; i < count; i++ {
		participants = append(participants, models.Participant{
			FullName:          randomFullName(),
			Email:             randomEmail(event.ID, i),
			DateOfBirth:       randomDOB(),
			SourceOfDiscovery: randomDiscoverySource(),
			EventID:           event.ID,
		})
	}

	const batchSize = 50
	return s.participantRepo.CreateInBatches(participants, batchSize)
}

var (
	seededRand       = rand.New(rand.NewSource(time.Now().UnixNano()))
	firstNames       = []string{"Alex", "Taylor", "Jordan", "Morgan", "Casey", "Riley", "Quinn", "Jamie", "Avery", "Parker"}
	lastNames        = []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Miller", "Davis", "Garcia", "Rodriguez", "Wilson"}
	discoverySources = []models.SourceOfDiscovery{
		models.SourceSocialMedia,
		models.SourceFriends,
		models.SourceFoundMyself,
	}
)

func randomFullName() string {
	first := firstNames[seededRand.Intn(len(firstNames))]
	last := lastNames[seededRand.Intn(len(lastNames))]
	return first + " " + last
}

func randomEmail(eventID string, index int) string {
	base := strings.ToLower(strings.ReplaceAll(randomFullName(), " ", "."))
	return fmt.Sprintf("%s+%s-%d@example.com", base, eventID[:8], index)
}

func randomDOB() *time.Time {
	year := seededRand.Intn(30) + 1975
	month := time.Month(seededRand.Intn(12) + 1)
	day := seededRand.Intn(28) + 1
	dob := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return &dob
}

func randomDiscoverySource() models.SourceOfDiscovery {
	return discoverySources[seededRand.Intn(len(discoverySources))]
}

type participantService struct {
	participantRepo repositories.ParticipantRepository
	eventRepo       repositories.EventRepository
}

// NewParticipantService creates a new participant service
func NewParticipantService(
	participantRepo repositories.ParticipantRepository,
	eventRepo repositories.EventRepository,
) ParticipantService {
	return &participantService{
		participantRepo: participantRepo,
		eventRepo:       eventRepo,
	}
}

func (s *participantService) RegisterParticipant(participant *models.Participant) (*models.Participant, error) {
	// Check if event exists
	_, err := s.eventRepo.FindByID(participant.EventID)
	if err != nil {
		return nil, err
	}

	// Set current time if DateOfBirth is not provided
	if participant.DateOfBirth == nil {
		now := time.Now()
		participant.DateOfBirth = &now
	}

	// Save participant
	if err := s.participantRepo.Create(participant); err != nil {
		return nil, err
	}

	return participant, nil
}

func (s *participantService) GetEventParticipants(eventID string) ([]*models.Participant, error) {
	return s.participantRepo.FindByEventID(eventID)
}

func (s *participantService) GetParticipantByID(id string) (*models.Participant, error) {
	return s.participantRepo.FindByID(id)
}

func (s *participantService) GetParticipantByEmail(email string) ([]*models.Participant, error) {
	return s.participantRepo.FindByEmail(email)
}

func (s *participantService) GetEventParticipantCount(eventID string) (int64, error) {
	return s.participantRepo.CountByEventID(eventID)
}

func (s *participantService) DeleteParticipant(id string) error {
	return s.participantRepo.Delete(id)
}
