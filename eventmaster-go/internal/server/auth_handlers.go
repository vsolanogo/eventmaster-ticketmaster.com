package server

import (
	"eventmaster-go/internal/models"
	"eventmaster-go/internal/services"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (s *Server) handleRegister(authService services.AuthService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req RegisterRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}

		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		user, err := authService.Register(req.Email, req.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusCreated, user.ToResponse())
	}
}

func (s *Server) handleLogin(authService services.AuthService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}

		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		session, err := authService.Login(req.Email, req.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		}

		cookie := &http.Cookie{
			Name:     s.config.SessionCookieName,
			Value:    session.Token,
			Path:     "/",
			Expires:  session.ExpiresAt,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, session.ToResponse(&session.User))
	}
}

func (s *Server) handleGetCurrentUser(authService services.AuthService) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, _ := c.Get("currentUser").(*models.User)
		if user == nil {
			userID, _ := c.Get("userID").(string)
			if userID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}

			var err error
			user, err = authService.GetUserByID(userID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
			}
		}

		return c.JSON(http.StatusOK, user.ToResponse())
	}
}

// handleLogout handles user logout
func (s *Server) handleLogout(authService services.AuthService) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(s.config.SessionCookieName)
		if err == nil && cookie.Value != "" {
			if err := authService.Logout(cookie.Value); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to logout")
			}
		}

		c.SetCookie(&http.Cookie{
			Name:     s.config.SessionCookieName,
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})

		return c.JSON(http.StatusOK, map[string]string{
			"message": "Successfully logged out",
		})
	}
}
