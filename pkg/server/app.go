// Package server implements a web server
package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"

	"github.com/erdnaxeli/wishlister"
)

// Server expose a single method Run() to run the web server.
type Server struct {
	logger     slog.Logger
	router     chi.Router
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

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	validate := validator.New(validator.WithRequiredStructEnabled())

	s := Server{
		logger:     *slog.New(slog.NewTextHandler(os.Stderr, nil)),
		router:     router,
		templates:  templates,
		validate:   validate,
		wishlister: config.Wishlister,
	}

	s.setRoutes()
	// s.setStatics()

	return s
}

// Run starts the server.
//
// It blocks forever.
func (s Server) Run() {
	s.logger.Error("error while running server", "err", http.ListenAndServe(":3000", s.router))
}
