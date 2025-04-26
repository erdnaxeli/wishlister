package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/erdnaxeli/wishlister"
)

type createWishListForm struct {
	Name  string `form:"name"  validate:"required,max=255"`
	User  string `form:"user"  validate:"required,max=255"`
	Email string `form:"email" validate:"max=255"`
}

func (s Server) createNewWishList(c echo.Context) error {
	form := createWishListForm{}
	err := c.Bind(&form)
	if err != nil {
		return err
	}

	err = s.validate.Struct(form)
	if err != nil {
		return err
	}

	params := wishlister.CreateWishlistParams{
		Name:      form.Name,
		Username:  form.User,
		UserEmail: form.Email,
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

type getWishListParam struct {
	AdminID string `param:"adminID"`
	ListID  string `param:"listID"`
}

func (s Server) getWishList(c echo.Context) error {
	params := getWishListParam{}
	err := c.Bind(&params)
	if err != nil {
		return err
	}

	var list wishlister.WishList
	ctx := c.Request().Context()

	if params.AdminID == "" {
		list, err = s.wishlister.GetWishList(ctx, params.ListID)
		if err != nil {
			return err
		}
	} else {
		list, err = s.wishlister.GetEditableWishList(ctx, params.ListID, params.AdminID)
		if err != nil {
			if errors.Is(err, wishlister.WishListInvalidAdminIDError{}) {
				return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/%s", params.ListID))
			}

			return err
		}
	}

	return renderOK(c, s.templates.RenderListViewBytes, list)
}
