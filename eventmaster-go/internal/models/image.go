package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Image represents an image in the system
type Image struct {
	Base
	Link string `json:"link" gorm:"not null"`
}

// ImageResponse represents the image data sent to clients
type ImageResponse struct {
	ID        string    `json:"id"`
	Link      string    `json:"link"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToResponse converts Image to ImageResponse
func (i *Image) ToResponse() *ImageResponse {
	return &ImageResponse{
		ID:        i.ID,
		Link:      i.Link,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}

// GenerateID generates a new UUID for the model
func GenerateID() string {
	return uuid.New().String()
}

// BeforeCreate is a hook that runs before creating an image
func (i *Image) BeforeCreate(tx *gorm.DB) error {
	i.ID = GenerateID()
	return nil
}
