package server

import (
	"io"
	"net/http"
)

func (s Server) setRoutes() {
	s.router.Get("/", s.renderOKFunc(s.templates.RenderIndex, nil))

	s.router.Get("/login", s.renderOKFunc(s.templates.RenderLogin, nil))
	s.router.Post("/login", s.sendMagicLink)
	s.router.Get("/login/magic/{token}", s.handleMagicLink)
	s.router.Get("/logout", s.logout)
	s.router.Get("/lists", s.getUserWishLists)

	s.router.Get("/new", s.getNewWishList)
	s.router.Post("/new", s.createNewWishList)

	s.router.Get("/group/new", s.renderOKFunc(s.templates.RenderNewGroup, nil))
	s.router.Post("/group/new", s.createNewGroup)

	s.router.Get("/l/{listID}", s.getWishList)
	s.router.Get("/l/{listID}/{adminID}", s.getWishList)
	s.router.Get("/l/{listID}/{adminID}/edit", s.editList)
	s.router.Post("/l/{listID}/{adminID}/edit", s.editList)

	// 404 page
	s.router.Get("/*", s.renderFunc(http.StatusNotFound, s.templates.RenderNotFoundError, nil))
}

func (s Server) renderOKFunc(templateFunc func(io.Writer, any) error, data any) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		s.renderOK(w, templateFunc, data)
	}
}

func (s Server) renderFunc(
	code int,
	templateFunc func(io.Writer, any) error,
	data any,
) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		s.render(w, code, templateFunc, data)
	}
}

func (s Server) renderOK(w http.ResponseWriter, templateFunc func(io.Writer, any) error, data any) {
	s.render(w, http.StatusOK, templateFunc, data)
}

func (s Server) render(
	w http.ResponseWriter,
	code int,
	templateFunc func(io.Writer, any) error,
	data any,
) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	err := templateFunc(w, data)
	if err != nil {
		s.logger.Error("error while rendering template", "err", err)
	}
}
