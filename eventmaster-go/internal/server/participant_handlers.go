package server

import (
	"eventmaster-go/internal/models"
	"eventmaster-go/internal/services"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// RegisterParticipantRequest represents the request body for registering a participant
type RegisterParticipantRequest struct {
	FullName          string                   `json:"fullName" validate:"required"`
	Email             string                   `json:"email" validate:"required,email"`
	DateOfBirth       *time.Time               `json:"dateOfBirth"`
	SourceOfDiscovery models.SourceOfDiscovery `json:"sourceOfDiscovery" validate:"required,oneof=social_media friends found_myself"`
}

// ParticipantResponse represents the participant response
type ParticipantResponse struct {
	ID                string                   `json:"id"`
	FullName          string                   `json:"fullName"`
	Email             string                   `json:"email"`
	DateOfBirth       *time.Time               `json:"dateOfBirth"`
	SourceOfDiscovery models.SourceOfDiscovery `json:"sourceOfDiscovery"`
	EventID           string                   `json:"eventId"`
	CreatedAt         time.Time                `json:"createdAt"`
	UpdatedAt         time.Time                `json:"updatedAt"`
}

// RegisterParticipantHandlers registers participant-related HTTP handlers
func (s *Server) RegisterParticipantHandlers(participantService services.ParticipantService) {
	participantGroup := s.apiGroup.Group("/participant")

	participantGroup.GET("/event/:eventId", s.handleGetEventParticipants(participantService))
	participantGroup.GET("/event/:eventId/registrations-per-day", s.handleRegistrationsPerDay(participantService))
	participantGroup.POST("/event/:eventId", s.handleRegisterParticipant(participantService))
	participantGroup.GET("/:id", s.handleGetParticipant(participantService))

	protected := participantGroup.Group("")
	protected.Use(s.requireAuth)
	{
		protected.DELETE("/:id", s.handleDeleteParticipant(participantService))
	}
}

func (s *Server) handleRegisterParticipant(svc services.ParticipantService) echo.HandlerFunc {
	return func(c echo.Context) error {
		eventID := c.Param("eventId")
		if eventID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "event ID is required")
		}

		var req RegisterParticipantRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}

		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		participant := &models.Participant{
			FullName:          req.FullName,
			Email:             req.Email,
			DateOfBirth:       req.DateOfBirth,
			SourceOfDiscovery: req.SourceOfDiscovery,
			EventID:           eventID,
		}

		savedParticipant, err := svc.RegisterParticipant(participant)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to register participant: "+err.Error())
		}

		return c.JSON(http.StatusCreated, savedParticipant.ToResponse())
	}
}

func (s *Server) handleGetEventParticipants(svc services.ParticipantService) echo.HandlerFunc {
	return func(c echo.Context) error {
		eventID := c.Param("eventId")
		if eventID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "event ID is required")
		}

		participants, err := svc.GetEventParticipants(eventID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch participants")
		}

		// Convert to response objects
		response := make([]*ParticipantResponse, len(participants))
		for i, p := range participants {
			resp := p.ToResponse()
			response[i] = &ParticipantResponse{
				ID:                resp.ID,
				FullName:          resp.FullName,
				Email:             resp.Email,
				DateOfBirth:       resp.DateOfBirth,
				SourceOfDiscovery: resp.SourceOfDiscovery,
				EventID:           resp.EventID,
				CreatedAt:         resp.CreatedAt,
				UpdatedAt:         resp.UpdatedAt,
			}
		}

		return c.JSON(http.StatusOK, response)
	}
}

func (s *Server) handleGetParticipant(svc services.ParticipantService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "participant ID is required")
		}

		participant, err := svc.GetParticipantByID(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "participant not found")
		}

		resp := participant.ToResponse()
		return c.JSON(http.StatusOK, resp)
	}
}

func (s *Server) handleDeleteParticipant(svc services.ParticipantService) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "participant ID is required")
		}

		if err := svc.DeleteParticipant(id); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete participant")
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (s *Server) handleRegistrationsPerDay(svc services.ParticipantService) echo.HandlerFunc {
	return func(c echo.Context) error {
		eventID := c.Param("eventId")
		if eventID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "event ID is required")
		}

		data, err := svc.RegistrationsPerDay(eventID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch registration analytics")
		}

		return c.JSON(http.StatusOK, data)
	}
}
