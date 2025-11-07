package repositories

import (
	"eventmaster-go/internal/models"

	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	BaseRepository[models.User]
	FindByEmail(email string) (*models.User, error)
	FindWithRoles(id string) (*models.User, error)
}

type userRepository struct {
	BaseRepository[models.User]
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	baseRepo := NewBaseRepository[models.User](db, models.User{})
	return &userRepository{
		BaseRepository: baseRepo,
		db:            db,
	}
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindWithRoles(id string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
