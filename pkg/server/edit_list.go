package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	nanoid "github.com/matoous/go-nanoid/v2"

	"github.com/erdnaxeli/wishlister"
)

type listEditTmplParams struct {
	Name string
	Data string
}

type editListData struct {
	Elements []editListDataElement `json:"elements"`
}

type editListDataElement struct {
	ID string `json:"id"`

	Name        string `json:"name"`
	NameError   string `json:"name_error"`
	Description string `json:"description"`
	URL         string `json:"url"`

	Error string `json:"error"`
}

// ErrInvalidForm is the error when the from sent is invalid, meaning expected data is
// not present. It probably means that the query was crafted and not sent through the
// HTML form.
var ErrInvalidForm = errors.New("invalid form")

func (s Server) editList(c echo.Context) error {
	params := getWishListParam{}
	err := c.Bind(&params)
	if err != nil {
		return err
	}

	list, err := s.wishlister.GetEditableWishList(
		c.Request().Context(),
		params.ListID,
		params.AdminID,
	)
	if err != nil {
		if errors.Is(err, wishlister.WishListNotFoundError{}) {
			return c.Render(http.StatusNotFound, "listNotFound", nil)
		}

		if errors.Is(err, wishlister.WishListInvalidAdminIDError{}) {
			return c.Render(http.StatusForbidden, "listAccessDenied", list)
		}

		return err
	}

	var data editListData

	if c.Request().Method == http.MethodPost {
		var ok bool
		data, ok, err = s.validateEditForm(c)
		if err != nil {
			return err
		}

		if ok {
			err := s.updateList(c, data, params.ListID, params.AdminID)
			if err != nil {
				return err
			}

			return c.Redirect(
				http.StatusSeeOther,
				fmt.Sprintf("/l/%s/%s", params.ListID, params.AdminID),
			)
		}
	} else {
		data = listToEditData(list)
	}

	dataJSON, err := json.Marshal(data.Elements)
	if err != nil {
		return err
	}

	tmplParams := listEditTmplParams{
		Name: list.Name,
		Data: string(dataJSON),
	}

	return renderOK(c, s.templates.RenderListEditBytes, tmplParams)
}

type listElementsForm struct {
	Descriptions []string `form:"description" validate:"required"`
	Names        []string `form:"name"        validate:"required"`
	Urls         []string `form:"url"         validate:"required"`
}

func (s Server) validateEditForm(c echo.Context) (editListData, bool, error) {
	data := editListData{}
	ok := true

	form := listElementsForm{}
	err := c.Bind(&form)
	if err != nil {
		return data, false, err
	}

	err = s.validate.Struct(form)
	if err != nil {
		return data, false, err
	}

	if len(form.Names) != len(form.Descriptions) || len(form.Names) != len(form.Urls) {
		return data, false, ErrInvalidForm
	}

	for i := range form.Names {
		element := editListDataElement{
			ID:          uuid.NewString(),
			Name:        form.Names[i],
			Description: form.Descriptions[i],
			URL:         form.Urls[i],
		}

		if form.Names[i] == "" {
			element.NameError = "Le nom ne peut pas être vide."
			ok = false
		}

		data.Elements = append(data.Elements, element)
	}

	return data, ok, nil
}

func (s Server) updateList(
	c echo.Context,
	form editListData,
	listID string,
	adminID string,
) error {
	elements := make([]wishlister.WishListElement, len(form.Elements))

	for idx, elt := range form.Elements {
		elements[idx] = wishlister.WishListElement{
			Name:        elt.Name,
			Description: elt.Description,
			URL:         elt.URL,
		}
	}

	return s.wishlister.UpdateListElements(c.Request().Context(), listID, adminID, elements)
}

func listToEditData(list wishlister.WishList) editListData {
	data := editListData{Elements: make([]editListDataElement, len(list.Elements))}
	for idx, element := range list.Elements {
		id, _ := nanoid.New()
		data.Elements[idx] = editListDataElement{
			ID:          id,
			Name:        element.Name,
			Description: element.Description,
			URL:         element.URL,
		}
	}

	return data
}
