package repositories

import (
	"gorm.io/gorm"
)

// BaseRepository defines common database operations
type BaseRepository[T any] interface {
	Create(entity *T) error
	FindByID(id string) (*T, error)
	Update(entity *T) error
	Delete(id string) error
}

type baseRepository[T any] struct {
	db    *gorm.DB
	model T
}

// NewBaseRepository creates a new base repository
func NewBaseRepository[T any](db *gorm.DB, model T) BaseRepository[T] {
	return &baseRepository[T]{
		db:    db,
		model: model,
	}
}

func (r *baseRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *baseRepository[T]) FindByID(id string) (*T, error) {
	var entity T
	err := r.db.First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *baseRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

func (r *baseRepository[T]) Delete(id string) error {
	return r.db.Delete(r.model, "id = ?", id).Error
}
