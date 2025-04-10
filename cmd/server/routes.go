package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/erdnaxeli/wishlister"
	"github.com/labstack/echo/v4"
)

func setRoutes(e *echo.Echo, app wishlister.App) {
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home", nil)
	})

	e.GET("/new", func(c echo.Context) error {
		return c.Render(http.StatusOK, "new", nil)
	})
	e.POST("/new", func(c echo.Context) error {
		return createNewWishList(c, app)
	})

	e.GET("/:listID", func(c echo.Context) error {
		return getWishList(c, app)
	})

	e.GET("/:listID/:adminID", func(c echo.Context) error {
		return getWishList(c, app)
	})

	e.GET("/:listID/:adminID/edit", func(c echo.Context) error {
		return editList(c, app)
	})

	e.POST("/:listID/:adminID/edit", func(c echo.Context) error {
		return editList(c, app)
	})
}

func createNewWishList(c echo.Context, app wishlister.App) error {
	listID, adminID, err := app.CreateWishList(
		c.Request().Context(),
		wishlister.CreateWishlistParams{
			Name:      c.FormValue("name"),
			Username:  c.FormValue("user"),
			UserEmail: c.FormValue("email"),
		},
	)
	if err != nil {
		return err
	}

	return c.Redirect(303, fmt.Sprintf("%s/%s", listID, adminID))
}

func getWishList(c echo.Context, app wishlister.App) error {
	listID := c.Param("listID")
	adminID := c.Param("adminID")

	var list wishlister.WishList
	var err error

	if adminID != "" {
		list, err = app.GetEditableWishList(c.Request().Context(), listID, adminID)
		if err != nil {
			if errors.Is(err, wishlister.WishListInvalidAdminIDError{}) {
				return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/%s", listID))
			}

			return err
		}
	} else {
		list, err = app.GetWishList(c.Request().Context(), listID)
		if err != nil {
			return err
		}
	}

	return c.Render(http.StatusOK, "listView", list)
}
