package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Session struct {
	Base
	UserID    string    `gorm:"not null"`
	Token     string    `gorm:"type:text;not null;unique"`
	ExpiresAt time.Time `gorm:"not null"`
	IP        string    `gorm:"size:45;not null"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for the Session model
func (Session) TableName() string {
	return "sessions"
}

// Claims represents the JWT claims structure
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// DTOs (Data Transfer Objects)
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SessionResponse struct {
	ID        string    `json:"id"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expires"`
}

func (s *Session) ToResponse() SessionResponse {
	return SessionResponse{
		ID:        s.ID,
		IP:        s.IP,
		CreatedAt: s.CreatedAt,
		ExpiresAt: s.ExpiresAt,
	}
}
