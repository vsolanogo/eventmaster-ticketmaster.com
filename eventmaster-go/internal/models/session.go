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

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      *UserResponse `json:"user"`
}

// ToResponse converts a Session to a LoginResponse
func (s *Session) ToResponse(user *User) *LoginResponse {
	return &LoginResponse{
		Token:     s.Token,
		ExpiresAt: s.ExpiresAt,
		User:      user.ToResponse(),
	}
}
