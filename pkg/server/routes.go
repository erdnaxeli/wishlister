package server

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) setRoutes() {
	s.e.GET("/", renderOKFunc(s.templates.RenderIndex, nil))

	s.e.GET("/login", renderOKFunc(s.templates.RenderLogin, nil))
	s.e.POST("/login", s.sendMagicLink)

	s.e.GET("/new", renderOKFunc(s.templates.RenderNew, nil))
	s.e.POST("/new", s.createNewWishList)

	s.e.GET("/group/new", renderOKFunc(s.templates.RenderNewGroup, nil))
	s.e.POST("/group/new", s.createNewGroup)

	s.e.GET("/l/:listID", s.getWishList)
	s.e.GET("/l/:listID/:adminID", s.getWishList)
	s.e.GET("/l/:listID/:adminID/edit", s.editList)
	s.e.POST("/l/:listID/:adminID/edit", s.editList)

	// 404 page
	s.e.GET("/*", renderFunc(http.StatusNotFound, s.templates.RenderNotFoundError, nil))
}

func renderOKFunc(templateFunc func(io.Writer, any) error, data any) func(echo.Context) error {
	return func(c echo.Context) error {
		return renderOK(c, templateFunc, data)
	}
}

func renderFunc(
	code int,
	templateFunc func(io.Writer, any) error,
	data any,
) func(echo.Context) error {
	return func(c echo.Context) error {
		return render(c, code, templateFunc, data)
	}
}

func renderOK(c echo.Context, templateFunc func(io.Writer, any) error, data any) error {
	return render(c, http.StatusOK, templateFunc, data)
}

func render(c echo.Context, code int, templateFunc func(io.Writer, any) error, data any) error {
	response := c.Response()

	header := response.Header()
	if header.Get(echo.HeaderContentType) == "" {
		header.Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	}

	response.WriteHeader(code)
	err := templateFunc(response.Writer, data)
	if err != nil {
		return err
	}

	return nil
}
