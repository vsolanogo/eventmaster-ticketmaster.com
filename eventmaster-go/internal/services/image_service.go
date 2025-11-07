package services

import (
	"eventmaster-go/internal/models"
	"eventmaster-go/internal/repositories"
)

// ImageService handles operations related to image records
// It complements FileService by allowing images to be registered using external links
// without downloading them locally.
type ImageService interface {
	CreateImagesWithLinks(links []string) ([]*models.Image, error)
}

type imageService struct {
	imageRepo repositories.ImageRepository
}

// NewImageService creates a new image service instance
func NewImageService(imageRepo repositories.ImageRepository) ImageService {
	return &imageService{
		imageRepo: imageRepo,
	}
}

func (s *imageService) CreateImagesWithLinks(links []string) ([]*models.Image, error) {
	created := make([]*models.Image, 0, len(links))

	for _, link := range links {
		if link == "" {
			continue
		}

		existing, err := s.imageRepo.FindByLink(link)
		if err == nil && existing != nil {
			created = append(created, existing)
			continue
		}

		image := &models.Image{Link: link}
		if err := s.imageRepo.Create(image); err != nil {
			// Skip this image but continue processing remaining links
			continue
		}
		created = append(created, image)
	}

	return created, nil
}
