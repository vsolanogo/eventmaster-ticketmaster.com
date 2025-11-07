package models

import (
	"time"

	"gorm.io/gorm"
)

// Base contains common fields for all models
type Base struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
