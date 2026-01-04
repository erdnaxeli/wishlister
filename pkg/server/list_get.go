package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/erdnaxeli/wishlister"
)

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
			if errors.Is(err, wishlister.ErrWishListInvalidAdminID) {
				return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/%s", params.ListID))
			}

			return err
		}
	}

	return renderOK(c, s.templates.RenderListView, list)
}
