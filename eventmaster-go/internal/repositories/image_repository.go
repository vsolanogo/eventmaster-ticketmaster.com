package repositories

import (
	"eventmaster-go/internal/models"
	"gorm.io/gorm"
)

// ImageRepository defines the interface for image data operations
type ImageRepository interface {
	BaseRepository[models.Image]
	FindByLink(link string) (*models.Image, error)
	AttachToEvent(event *models.Event, images []*models.Image) error
}

type imageRepository struct {
	BaseRepository[models.Image]
	db *gorm.DB
}

// NewImageRepository creates a new image repository
func NewImageRepository(db *gorm.DB) ImageRepository {
	baseRepo := NewBaseRepository[models.Image](db, models.Image{})
	return &imageRepository{
		BaseRepository: baseRepo,
		db:            db,
	}
}

func (r *imageRepository) FindByLink(link string) (*models.Image, error) {
	var image models.Image
	err := r.db.Where("link = ?", link).First(&image).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &image, nil
}

func (r *imageRepository) AttachToEvent(event *models.Event, images []*models.Image) error {
	if event == nil {
		return nil
	}
	return r.db.Model(event).Association("Images").Append(images)
}
