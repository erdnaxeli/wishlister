package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/erdnaxeli/wishlister"
)

type createGroupForm struct {
	Name  string `form:"name"  validate:"required,max=255"`
	Email string `form:"email" validate:"required,email,max=255"`
}

func (s Server) createNewGroup(c echo.Context) error {
	form := createGroupForm{}
	err := c.Bind(&form)
	if err != nil {
		return err
	}

	err = s.validate.Struct(form)
	if err != nil {
		return err
	}

	groupID, err := s.wishlister.CreateGroup(
		c.Request().Context(),
		wishlister.CreateGroupParams{
			Name:      form.Name,
			UserEmail: form.Email,
		},
	)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/g/%s", groupID))
}
