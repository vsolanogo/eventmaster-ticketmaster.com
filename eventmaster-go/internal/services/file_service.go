package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"eventmaster-go/internal/models"
	"eventmaster-go/internal/repositories"
	"github.com/google/uuid"
)

// FileService handles file uploads and management
type FileService interface {
	SaveUploadedFile(file *multipart.FileHeader) (*models.Image, error)
	SaveFileFromURL(url string) (*models.Image, error)
	GetImageByID(id string) (*models.Image, error)
	DeleteImage(id string) error
	ResolveFilePath(link string) string
}

func (s *fileService) ResolveFilePath(link string) string {
	trimmed := strings.TrimPrefix(link, s.baseURL)
	trimmed = strings.TrimLeft(trimmed, "/")
	return filepath.Join(s.uploadPath, trimmed)
}

type fileService struct {
	imageRepo    repositories.ImageRepository
	uploadPath   string
	baseURL      string
	allowedTypes map[string]bool
}

// NewFileService creates a new file service
func NewFileService(
	imageRepo repositories.ImageRepository,
	uploadPath string,
	baseURL string,
) FileService {
	// Ensure upload directory exists
	os.MkdirAll(uploadPath, 0755)

	return &fileService{
		imageRepo:  imageRepo,
		uploadPath: uploadPath,
		baseURL:    strings.TrimRight(baseURL, "/"),
		allowedTypes: map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
			"image/webp": true,
		},
	}
}

func (s *fileService) SaveUploadedFile(file *multipart.FileHeader) (*models.Image, error) {
	// Validate file type
	if !s.allowedTypes[file.Header.Get("Content-Type")] {
		return nil, fmt.Errorf("file type not allowed")
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Generate a unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(s.uploadPath, filename)

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	// Copy the file content
	if _, err = io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	publicLink := fmt.Sprintf("%s/%s", strings.TrimRight(s.baseURL, "/"), filename)
	if !strings.HasPrefix(publicLink, "/") {
		publicLink = "/" + publicLink
	}

	// Create image record in database
	image := &models.Image{
		Link: publicLink,
	}

	if err := s.imageRepo.Create(image); err != nil {
		// Clean up the uploaded file if database operation fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save image record: %v", err)
	}

	return image, nil
}

func (s *fileService) SaveFileFromURL(url string) (*models.Image, error) {
	// Generate a unique filename based on URL hash and timestamp
	hash := md5.Sum([]byte(url + time.Now().String()))
	ext := filepath.Ext(url)
	if ext == "" {
		ext = ".jpg" // Default extension if none provided
	}
	filename := hex.EncodeToString(hash[:]) + ext
	filePath := filepath.Join(s.uploadPath, filename)

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	if !s.allowedTypes[contentType] {
		return nil, fmt.Errorf("file type not allowed: %s", contentType)
	}

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	// Copy the file content
	if _, err = io.Copy(dst, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	publicLink := fmt.Sprintf("%s/%s", strings.TrimRight(s.baseURL, "/"), filename)
	if !strings.HasPrefix(publicLink, "/") {
		publicLink = "/" + publicLink
	}

	// Create image record in database
	image := &models.Image{
		Link: publicLink,
	}

	if err := s.imageRepo.Create(image); err != nil {
		// Clean up the downloaded file if database operation fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save image record: %v", err)
	}

	return image, nil
}

func (s *fileService) GetImageByID(id string) (*models.Image, error) {
	return s.imageRepo.FindByID(id)
}

func (s *fileService) DeleteImage(id string) error {
	image, err := s.imageRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Delete the file
	filename := filepath.Base(image.Link)
	filePath := filepath.Join(s.uploadPath, filename)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete image file: %v", err)
	}

	// Delete the database record
	return s.imageRepo.Delete(id)
}
