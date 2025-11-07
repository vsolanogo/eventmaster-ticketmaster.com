package models

import (
	"time"

	"gorm.io/gorm"
)

// SourceOfDiscovery represents how a participant discovered the event
type SourceOfDiscovery string

const (
	SourceSocialMedia SourceOfDiscovery = "social_media"
	SourceFriends    SourceOfDiscovery = "friends"
	SourceFoundMyself SourceOfDiscovery = "found_myself"
)

// Participant represents an event participant
type Participant struct {
	Base
	FullName           string          `json:"fullName" gorm:"not null"`
	Email              string          `json:"email" gorm:"not null;index"`
	DateOfBirth        *time.Time      `json:"dateOfBirth" gorm:"not null"`
	SourceOfDiscovery  SourceOfDiscovery `json:"sourceOfDiscovery" gorm:"type:varchar(50);not null"`
	EventID            string          `json:"eventId" gorm:"type:uuid;not null;index"`
	Event              *Event          `json:"-" gorm:"foreignKey:EventID"`
}

// ParticipantResponse represents the participant data sent to clients
type ParticipantResponse struct {
	ID                string          `json:"id"`
	FullName          string          `json:"fullName"`
	Email             string          `json:"email"`
	DateOfBirth       *time.Time      `json:"dateOfBirth"`
	SourceOfDiscovery SourceOfDiscovery `json:"sourceOfDiscovery"`
	EventID           string          `json:"eventId"`
	CreatedAt         time.Time       `json:"createdAt"`
	UpdatedAt         time.Time       `json:"updatedAt"`
}

// ToResponse converts Participant to ParticipantResponse
func (p *Participant) ToResponse() *ParticipantResponse {
	return &ParticipantResponse{
		ID:                p.ID,
		FullName:          p.FullName,
		Email:             p.Email,
		DateOfBirth:       p.DateOfBirth,
		SourceOfDiscovery: p.SourceOfDiscovery,
		EventID:           p.EventID,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
}

// BeforeCreate is a hook that runs before creating a participant
func (p *Participant) BeforeCreate(tx *gorm.DB) error {
	p.ID = GenerateID()
	return nil
}
