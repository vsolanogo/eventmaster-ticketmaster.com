package repositories

import (
	"eventmaster-go/internal/models"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EventRepository defines the interface for event data operations
type EventRepository interface {
	BaseRepository[models.Event]
	FindByDateRange(start, end time.Time) ([]*models.Event, error)
	FindByUserID(userID string) ([]*models.Event, error)
	FindWithImages(id string) (*models.Event, error)
	FindByExternalID(externalID string) (*models.Event, error)
	FindPaginated(page, limit int, sortBy, sortOrder string) ([]*models.Event, int64, error)
}

type eventRepository struct {
	BaseRepository[models.Event]
	db *gorm.DB
}

// NewEventRepository creates a new event repository
func NewEventRepository(db *gorm.DB) EventRepository {
	baseRepo := NewBaseRepository[models.Event](db, models.Event{})
	return &eventRepository{
		BaseRepository: baseRepo,
		db:             db,
	}
}

func (r *eventRepository) FindByDateRange(start, end time.Time) ([]*models.Event, error) {
	var events []*models.Event
	err := r.db.Where("event_date BETWEEN ? AND ?", start, end).
		Preload("Images").
		Preload("User").
		Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventRepository) FindByUserID(userID string) ([]*models.Event, error) {
	var events []*models.Event
	err := r.db.Where("user_id = ?", userID).
		Preload("Images").
		Preload("User").
		Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventRepository) FindWithImages(id string) (*models.Event, error) {
	var event models.Event
	err := r.db.Preload("Images").
		Preload("User").
		First(&event, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// FindByExternalID finds an event by its external ID
func (r *eventRepository) FindByExternalID(externalID string) (*models.Event, error) {
	var event models.Event
	err := r.db.First(&event, "external_id = ?", externalID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) FindPaginated(page, limit int, sortBy, sortOrder string) ([]*models.Event, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var total int64
	if err := r.db.Model(&models.Event{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortColumnMap := map[string]string{
		"eventDate": "event_date",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
		"title":     "title",
	}

	columnName := "event_date"
	if mapped, ok := sortColumnMap[sortBy]; ok {
		columnName = mapped
	} else if sortBy != "" {
		columnName = sortBy
	}

	orderClause := clause.OrderByColumn{
		Column: clause.Column{Name: columnName},
		Desc:   strings.EqualFold(sortOrder, "DESC"),
	}

	var events []*models.Event
	err := r.db.Preload("Images").
		Preload("User").
		Order(orderClause).
		Offset(offset).
		Limit(limit).
		Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}
