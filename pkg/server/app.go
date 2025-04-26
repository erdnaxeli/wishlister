// Package server implements a web server
package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/erdnaxeli/wishlister"
)

// Server expose a single method Run() to run the web server.
type Server struct {
	e          *echo.Echo
	templates  Templates
	validate   *validator.Validate
	wishlister wishlister.App
}

// Config is the server configuration.
type Config struct {
	Wishlister wishlister.App
}

// New creates a new Server object.
func New(config Config) Server {
	templates := NewTemplates()

	e := echo.New()
	e.Debug = true
	e.Pre(
		middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{RedirectCode: 308}),
	)
	e.Use(middleware.Logger())

	validate := validator.New(validator.WithRequiredStructEnabled())

	s := Server{
		e:          e,
		templates:  templates,
		validate:   validate,
		wishlister: config.Wishlister,
	}

	s.setRoutes()
	s.setStatics()

	return s
}

// Run starts the server.
//
// It blocks forever.
func (s Server) Run() {
	s.e.Logger.Fatal(s.e.Start(":8080"))
}
