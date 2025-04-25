package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/erdnaxeli/wishlister"
)

func setRoutes(e *echo.Echo, app wishlister.App, templates Templates) {
	e.GET("/", renderOKFunc(templates.RenderIndexBytes, nil))
	e.GET("/new", renderOKFunc(templates.RenderNewBytes, nil))
	e.GET("/group/new", renderOKFunc(templates.RenderNewGroupBytes, nil))

	e.POST("/new", func(c echo.Context) error {
		return createNewWishList(c, app)
	})
	e.POST("/group/new", func(c echo.Context) error {
		return createNewGroup(c, app)
	})

	e.GET("/l/:listID", func(c echo.Context) error {
		return getWishList(c, app, templates)
	})

	e.GET("/l/:listID/:adminID", func(c echo.Context) error {
		return getWishList(c, app, templates)
	})

	e.GET("/l/:listID/:adminID/edit", func(c echo.Context) error {
		return editList(c, app, templates)
	})

	e.POST("/l/:listID/:adminID/edit", func(c echo.Context) error {
		return editList(c, app, templates)
	})

	e.GET("/*", renderFunc(http.StatusNotFound, templates.RenderNotFoundErrorBytes, nil))
}

func createNewWishList(c echo.Context, app wishlister.App) error {
	params := wishlister.CreateWishlistParams{
		Name:      c.FormValue("name"),
		Username:  c.FormValue("user"),
		UserEmail: c.FormValue("email"),
	}

	listID, adminID, err := app.CreateWishList(
		c.Request().Context(),
		params,
	)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("l/%s/%s", listID, adminID))
}

func getWishList(c echo.Context, app wishlister.App, templates Templates) error {
	listID := c.Param("listID")
	adminID := c.Param("adminID")

	var list wishlister.WishList
	var err error

	if adminID == "" {
		list, err = app.GetWishList(c.Request().Context(), listID)
		if err != nil {
			return err
		}
	} else {
		list, err = app.GetEditableWishList(c.Request().Context(), listID, adminID)
		if err != nil {
			if errors.Is(err, wishlister.WishListInvalidAdminIDError{}) {
				return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/%s", listID))
			}

			return err
		}
	}

	return renderOK(c, templates.RenderListViewBytes, list)
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
