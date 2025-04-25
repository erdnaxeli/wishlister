package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) setRoutes() {
	s.e.GET("/", renderOKFunc(s.templates.RenderIndexBytes, nil))
	s.e.GET("/new", renderOKFunc(s.templates.RenderNewBytes, nil))
	s.e.GET("/group/new", renderOKFunc(s.templates.RenderNewGroupBytes, nil))

	s.e.POST("/new", s.createNewWishList)
	s.e.POST("/group/new", s.createNewGroup)

	s.e.GET("/l/:listID", s.getWishList)
	s.e.GET("/l/:listID/:adminID", s.getWishList)
	s.e.GET("/l/:listID/:adminID/edit", s.editList)
	s.e.POST("/l/:listID/:adminID/edit", s.editList)

	// 404 page
	s.e.GET("/*", renderFunc(http.StatusNotFound, s.templates.RenderNotFoundErrorBytes, nil))
}

func renderOKFunc(templateFunc func(data any) ([]byte, error), data any) func(echo.Context) error {
	return func(c echo.Context) error {
		return renderOK(c, templateFunc, data)
	}
}

func renderFunc(
	code int,
	templateFunc func(data any) ([]byte, error),
	data any,
) func(echo.Context) error {
	return func(c echo.Context) error {
		return render(c, code, templateFunc, data)
	}
}

func renderOK(c echo.Context, templateFunc func(data any) ([]byte, error), data any) error {
	return render(c, http.StatusOK, templateFunc, data)
}

func render(c echo.Context, code int, templateFunc func(data any) ([]byte, error), data any) error {
	bytes, err := templateFunc(data)
	if err != nil {
		return err
	}

	return c.HTMLBlob(code, bytes)
}
