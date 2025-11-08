package repositories

import (
	"eventmaster-go/internal/models"
	"gorm.io/gorm"
	"time"
)

// ParticipantRepository defines the interface for participant data operations
type ParticipantRepository interface {
	BaseRepository[models.Participant]
	FindByEventID(eventID string) ([]*models.Participant, error)
	FindByEmail(email string) ([]*models.Participant, error)
	CountByEventID(eventID string) (int64, error)
	CreateInBatches(participants []models.Participant, batchSize int) error
	RegistrationsPerDay(eventID string) ([]*RegistrationsPerDayResult, error)
}

type participantRepository struct {
	BaseRepository[models.Participant]
	db *gorm.DB
}

// NewParticipantRepository creates a new participant repository
func NewParticipantRepository(db *gorm.DB) ParticipantRepository {
	baseRepo := NewBaseRepository[models.Participant](db, models.Participant{})
	return &participantRepository{
		BaseRepository: baseRepo,
		db:             db,
	}
}

func (r *participantRepository) FindByEventID(eventID string) ([]*models.Participant, error) {
	var participants []*models.Participant
	err := r.db.Where("event_id = ?", eventID).Find(&participants).Error
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (r *participantRepository) FindByEmail(email string) ([]*models.Participant, error) {
	var participants []*models.Participant
	err := r.db.Where("email = ?", email).Find(&participants).Error
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (r *participantRepository) CountByEventID(eventID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Participant{}).Where("event_id = ?", eventID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *participantRepository) CreateInBatches(participants []models.Participant, batchSize int) error {
	if len(participants) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = len(participants)
	}
	return r.db.CreateInBatches(participants, batchSize).Error
}

type RegistrationsPerDayResult struct {
	Date  time.Time `json:"date"`
	Count int64     `json:"count"`
}

func (r *participantRepository) RegistrationsPerDay(eventID string) ([]*RegistrationsPerDayResult, error) {
	var results []*RegistrationsPerDayResult
	err := r.db.Model(&models.Participant{}).
		Select("DATE(created_at) AS date, COUNT(*) AS count").
		Where("event_id = ?", eventID).
		Group("DATE(created_at)").
		Order("date").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
