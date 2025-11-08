package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Base
	Email    string  `gorm:"size:100;not null;unique" validate:"required,email"`
	Password string  `gorm:"size:255;not null" validate:"required,min=8"`
	Roles    []*Role `gorm:"many2many:user_roles;"`
	Sessions []*Session
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate is a hook that runs before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Generate a hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the hashed password
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// DTOs (Data Transfer Objects)
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserResponse struct {
	ID        string            `json:"id"`
	Email     string            `json:"email"`
	Roles     []RoleResponse    `json:"role"`
	Sessions  []SessionResponse `json:"session"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

// ToResponse converts a User to a UserResponse
func (u *User) ToResponse() *UserResponse {
	roles := make([]RoleResponse, 0, len(u.Roles))
	for _, role := range u.Roles {
		if role == nil {
			continue
		}
		roles = append(roles, role.ToResponse())
	}

	sessions := make([]SessionResponse, 0, len(u.Sessions))
	for _, session := range u.Sessions {
		if session == nil {
			continue
		}
		sessions = append(sessions, session.ToResponse())
	}

	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Roles:     roles,
		Sessions:  sessions,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
