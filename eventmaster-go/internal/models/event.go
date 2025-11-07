package models

import (
	"time"

	"gorm.io/gorm"
)

// Event represents an event in the system
type Event struct {
	Base
	Title         string     `json:"title" gorm:"not null"`
	Description   string     `json:"description" gorm:"type:text"`
	Organizer     string     `json:"organizer" gorm:"not null"`
	EventDate     *time.Time `json:"eventDate" gorm:"not null"`
	Latitude      float64    `json:"latitude" gorm:"type:decimal(10,8)"`
	Longitude     float64    `json:"longitude" gorm:"type:decimal(11,8)"`
	UserID        string     `json:"userId" gorm:"type:uuid;not null"`
	User          User       `json:"-" gorm:"foreignKey:UserID"`
	Images        []Image    `json:"images" gorm:"many2many:event_images;"`
	Location      string     `json:"location" gorm:"type:text"`
	ExternalID    string     `json:"externalId" gorm:"index"`
	ExternalURL   string     `json:"externalUrl" gorm:"type:text"`
	EventType     string     `json:"eventType"`
	IsExternal    bool       `json:"isExternal" gorm:"default:false"`
}

// EventResponse represents the event data sent to clients
type EventResponse struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Organizer   string         `json:"organizer"`
	EventDate   *time.Time     `json:"eventDate"`
	Latitude    float64        `json:"latitude"`
	Longitude   float64        `json:"longitude"`
	Location    string         `json:"location,omitempty"`
	ExternalID  string         `json:"externalId,omitempty"`
	ExternalURL string         `json:"externalUrl,omitempty"`
	EventType   string         `json:"eventType,omitempty"`
	IsExternal  bool           `json:"isExternal"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Images      []ImageResponse `json:"images,omitempty"`
	User        *UserResponse  `json:"user,omitempty"`
}

// ToResponse converts Event to EventResponse
func (e *Event) ToResponse() *EventResponse {
	images := make([]ImageResponse, len(e.Images))
	for i, img := range e.Images {
		images[i] = *img.ToResponse()
	}

	var userResp *UserResponse
	if e.User.ID != "" {
		userResp = e.User.ToResponse()
	}

	return &EventResponse{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Organizer:   e.Organizer,
		EventDate:   e.EventDate,
		Latitude:    e.Latitude,
		Longitude:   e.Longitude,
		Location:    e.Location,
		ExternalID:  e.ExternalID,
		ExternalURL: e.ExternalURL,
		EventType:   e.EventType,
		IsExternal:  e.IsExternal,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		Images:      images,
		User:        userResp,
	}
}

// BeforeCreate is a hook that runs before creating an event
func (e *Event) BeforeCreate(tx *gorm.DB) error {
	e.ID = GenerateID()
	return nil
}
