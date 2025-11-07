package server

import (
	"eventmaster-go/internal/models"
	"eventmaster-go/internal/services"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// CreateEventRequest represents the request body for creating an event
type CreateEventRequest struct {
	Title       string     `json:"title" validate:"required,min=2,max=255"`
	Description string     `json:"description" validate:"omitempty,max=5000"`
	Organizer   string     `json:"organizer" validate:"omitempty,max=255"`
	EventDate   *time.Time `json:"eventDate" validate:"required"`
	Latitude    float64    `json:"latitude" validate:"required,gte=-90,lte=90"`
	Longitude   float64    `json:"longitude" validate:"required,gte=-180,lte=180"`
	ImageIDs    []string   `json:"imageIds" validate:"omitempty,dive,uuid4"`
}

// UpdateEventRequest represents the request body for updating an event
type UpdateEventRequest struct {
	Title       *string     `json:"title" validate:"omitempty,min=2,max=255"`
	Description *string     `json:"description" validate:"omitempty,max=5000"`
	Organizer   *string     `json:"organizer" validate:"omitempty,max=255"`
	EventDate   *time.Time  `json:"eventDate" validate:"omitempty"`
	Latitude    *float64    `json:"latitude" validate:"omitempty,gte=-90,lte=90"`
	Longitude   *float64    `json:"longitude" validate:"omitempty,gte=-180,lte=180"`
	ImageIDs    []string    `json:"imageIds" validate:"omitempty,dive,uuid4"`
}

// EventResponse represents the event response
type EventResponse struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Organizer   string            `json:"organizer"`
	EventDate   *time.Time        `json:"eventDate"`
	Latitude    float64           `json:"latitude"`
	Longitude   float64           `json:"longitude"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	Images      []*models.ImageResponse `json:"images,omitempty"`
}

// RegisterEventHandlers registers event-related HTTP handlers
func (s *Server) RegisterEventHandlers(eventService services.EventService) {
	eventGroup := s.apiGroup.Group("/events")
	
	// Public routes
	eventGroup.GET("", s.handleGetEvents(eventService))
	eventGroup.GET("/:id", s.handleGetEvent(eventService))

	// Protected routes (require authentication)
	protected := eventGroup.Group("")
	protected.Use(s.requireAuth)
	{
		protected.POST("", s.handleCreateEvent(eventService))
		protected.PUT("/:id", s.handleUpdateEvent(eventService))
		protected.DELETE("/:id", s.handleDeleteEvent(eventService))
	}
}

func (s *Server) handleCreateEvent(svc services.EventService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, _ := c.Get("userID").(string)
		if userID == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
		}

		var req CreateEventRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}

		if err := c.Validate(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		event := &models.Event{
			Title:       req.Title,
			Description: req.Description,
			Organizer:   req.Organizer,
			EventDate:   req.EventDate,
			Latitude:    req.Latitude,
			Longitude:   req.Longitude,
		}

		createdEvent, err := svc.CreateEvent(event, userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create event: "+err.Error())
		}

		return c.JSON(http.StatusCreated, createdEvent.ToResponse())
	}
}

func (s *Server) handleGetEvent(svc services.EventService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "event ID is required")
		}

		event, err := svc.GetEventByID(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "event not found")
		}

		return c.JSON(http.StatusOK, event.ToResponse())
	}
}

func (s *Server) handleGetEvents(svc services.EventService) echo.HandlerFunc {
	return func(c echo.Context) error {
		page := parseQueryInt(c, "page", 1)
		limit := parseQueryInt(c, "limit", 10)
		sortBy := c.QueryParam("sortBy")
		if sortBy == "" {
			sortBy = "event_date"
		}
		sortOrder := c.QueryParam("sortOrder")
		if sortOrder == "" {
			sortOrder = "ASC"
		}

		events, totalCount, err := svc.GetPaginatedEvents(page, limit, sortBy, sortOrder)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch events")
		}

		responses := make([]*models.EventResponse, len(events))
		for i, event := range events {
			responses[i] = event.ToResponse()
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"events":     responses,
			"totalCount": totalCount,
		})
	}
}

func parseQueryInt(c echo.Context, name string, defaultValue int) int {
	value := c.QueryParam(name)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 || parsed > math.MaxInt32 {
		return defaultValue
	}

	return parsed
}

func (s *Server) handleUpdateEvent(svc services.EventService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, _ := c.Get("userID").(string)
		if userID == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
		}

		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "event ID is required")
		}

		// Verify the event exists and belongs to the user
		event, err := svc.GetEventByID(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "event not found")
		}

		if event.UserID != userID {
			return echo.NewHTTPError(http.StatusForbidden, "not authorized to update this event")
		}

		var req UpdateEventRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}

		if err := c.Validate(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// Update only the fields that are provided in the request
		if req.Title != nil {
			event.Title = *req.Title
		}
		if req.Description != nil {
			event.Description = *req.Description
		}
		if req.Organizer != nil {
			event.Organizer = *req.Organizer
		}
		if req.EventDate != nil {
			event.EventDate = req.EventDate
		}
		if req.Latitude != nil {
			event.Latitude = *req.Latitude
		}
		if req.Longitude != nil {
			event.Longitude = *req.Longitude
		}
		if req.ImageIDs != nil {
			// TODO: handle image association updates similar to NestJS if needed
		}

		updatedEvent, err := svc.UpdateEvent(id, event)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update event")
		}

		return c.JSON(http.StatusOK, updatedEvent.ToResponse())
	}
}

func (s *Server) handleDeleteEvent(svc services.EventService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, _ := c.Get("userID").(string)
		if userID == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
		}

		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "event ID is required")
		}

		// Verify the event exists and belongs to the user
		event, err := svc.GetEventByID(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "event not found")
		}

		if event.UserID != userID {
			return echo.NewHTTPError(http.StatusForbidden, "not authorized to delete this event")
		}

		if err := svc.DeleteEvent(id); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete event")
		}

		return c.NoContent(http.StatusNoContent)
	}
}
