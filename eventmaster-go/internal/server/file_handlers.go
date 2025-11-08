package server

import (
	"eventmaster-go/internal/services"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

// RegisterFileHandlers registers file-related HTTP handlers
// Currently, file uploads are handled by the image controller in NestJS
// with routes like POST /images/upload and POST /images/upload/url
// This is a placeholder for future implementation
func (s *Server) RegisterFileHandlers(fileService services.FileService) {
	imageGroup := s.apiGroup.Group("/image")
	imageGroup.Use(s.requireAuth)

	imageGroup.POST("", s.handleFileUpload(fileService))
	imageGroup.POST("/url", s.handleFileUploadFromURL(fileService))
	imageGroup.DELETE("/:id", s.handleDeleteFile(fileService))
	imageGroup.GET("/:id", s.handleGetFile(fileService))
}

func (s *Server) handleFileUpload(svc services.FileService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the file from the form data
		file, err := c.FormFile("file")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "file is required")
		}

		// Save the file
		image, err := svc.SaveUploadedFile(file)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to save file: "+err.Error())
		}

		return c.JSON(http.StatusCreated, image.ToResponse())
	}
}

func (s *Server) handleFileUploadFromURL(svc services.FileService) echo.HandlerFunc {
	type request struct {
		URL string `json:"url" validate:"required,url"`
	}

	return func(c echo.Context) error {
		var req request
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}

		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// Download and save the file
		image, err := svc.SaveFileFromURL(req.URL)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to save file from URL: "+err.Error())
		}

		return c.JSON(http.StatusCreated, image.ToResponse())
	}
}

func (s *Server) handleGetFile(svc services.FileService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "file ID is required")
		}

		image, err := svc.GetImageByID(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "file not found")
		}

		// Check if the file exists
		filePath := filepath.Join(os.Getenv("UPLOAD_PATH"), filepath.Base(image.Link))
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return echo.NewHTTPError(http.StatusNotFound, "file not found on disk")
		}

		// Return the file
		return c.File(filePath)
	}
}

func (s *Server) handleDeleteFile(svc services.FileService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("userID").(string)
		if userID == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
		}

		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "file ID is required")
		}

		// Check if the file is used in any event before deleting
		// This would require adding a method to check for image usage
		// For now, we'll just delete the file

		if err := svc.DeleteImage(id); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete file: "+err.Error())
		}

		return c.NoContent(http.StatusNoContent)
	}
}
