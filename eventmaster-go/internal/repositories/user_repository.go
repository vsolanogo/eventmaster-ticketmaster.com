package repositories

import (
	"errors"

	"eventmaster-go/internal/models"

	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	BaseRepository[models.User]
	FindByEmail(email string) (*models.User, error)
	FindWithAssociations(id string) (*models.User, error)
	AttachRoleByName(user *models.User, roleName string) error
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
		db:             db,
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

func (r *userRepository) FindWithAssociations(id string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Preload("Sessions").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) AttachRoleByName(user *models.User, roleName string) error {
	if user == nil || user.ID == "" {
		return errors.New("user must have an ID before attaching roles")
	}

	var role models.Role
	err := r.db.Where("name = ?", roleName).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			role = models.Role{Name: roleName}
			if err := r.db.Create(&role).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return r.db.Model(user).Association("Roles").Append(&role)
}
