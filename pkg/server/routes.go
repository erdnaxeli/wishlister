package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/erdnaxeli/wishlister"
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

func (s Server) createNewWishList(c echo.Context) error {
	params := wishlister.CreateWishlistParams{
		Name:      c.FormValue("name"),
		Username:  c.FormValue("user"),
		UserEmail: c.FormValue("email"),
	}

	listID, adminID, err := s.wishlister.CreateWishList(
		c.Request().Context(),
		params,
	)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("l/%s/%s", listID, adminID))
}

func (s Server) getWishList(c echo.Context) error {
	listID := c.Param("listID")
	adminID := c.Param("adminID")

	var list wishlister.WishList
	var err error

	if adminID == "" {
		list, err = s.wishlister.GetWishList(c.Request().Context(), listID)
		if err != nil {
			return err
		}
	} else {
		list, err = s.wishlister.GetEditableWishList(c.Request().Context(), listID, adminID)
		if err != nil {
			if errors.Is(err, wishlister.WishListInvalidAdminIDError{}) {
				return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/%s", listID))
			}

			return err
		}
	}

	return renderOK(c, s.templates.RenderListViewBytes, list)
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
