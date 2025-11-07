package server

import (
	"context"
	"eventmaster-go/internal/services"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo        *echo.Echo
	apiGroup    *echo.Group
	config      Config
	authService services.AuthService
}

type Config struct {
	Port              string
	SessionCookieName string
}

func NewServer(authService services.AuthService, config Config) *Server {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	e.Validator = &CustomValidator{Validator: NewValidator()}

	apiGroup := e.Group("/api")

	server := &Server{
		echo:        e,
		apiGroup:    apiGroup,
		config:      config,
		authService: authService,
	}

	// Setup routes
	server.setupRoutes(authService)

	return server
}

func (s *Server) setupRoutes(authService services.AuthService) {
	// Health check
	s.apiGroup.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Auth routes
	s.apiGroup.POST("/register", s.handleRegister(authService))
	s.apiGroup.POST("/login", s.handleLogin(authService))
	s.apiGroup.POST("/logout", s.requireAuth(s.handleLogout(authService)))

	// Protected user route
	s.apiGroup.GET("/user", s.requireAuth(s.handleGetCurrentUser(authService)))
}

func (s *Server) requireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(s.config.SessionCookieName)
		if err != nil || cookie.Value == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}

		user, err := s.authService.ValidateSession(cookie.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}

		c.Set("userID", user.ID)
		c.Set("currentUser", user)

		return next(c)
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.echo.Start(fmt.Sprintf(":%s", s.config.Port))
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
