package server

import (
	"errors"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/go-playground/form"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	nanoid "github.com/matoous/go-nanoid/v2"

	"github.com/erdnaxeli/wishlister"
	"github.com/erdnaxeli/wishlister/pkg/server/templates"
)

// ErrInvalidForm is the error when the form sent is invalid, meaning expected data is
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
			return render(c, http.StatusNotFound, s.templates.RenderListNotFoundBytes, nil)
		}

		if errors.Is(err, wishlister.WishListInvalidAdminIDError{}) {
			return render(c, http.StatusForbidden, s.templates.RenderListAccessDeniedBytes, nil)
		}

		return err
	}

	var data templates.ListEditForm

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

	data.Name = list.Name
	s.renderer.Render()
	return renderOK(c, s.renderer.RenderListEditBytes, data)
}

func (s Server) validateEditForm(c echo.Context) (templates.ListEditForm, bool, error) {
	data := templates.ListEditForm{}
	ok := true
	decoder := form.NewDecoder()
	values, err := c.FormParams()
	if err != nil {
		c.Logger().Print("Error while reading wishlist elements form: %s", err)
		return data, false, err
	}

	err = decoder.Decode(&data, values)
	if err != nil {
		c.Logger().Print("Error while decoding wishlist elements form: %s", err)
		return data, false, err
	}

	for i, element := range data.Elements {
		element.ID = uuid.NewString()
		element, ok = s.validateElement(element, ok)
		data.Elements[i] = element
	}

	return data, ok, nil
}

func (s Server) updateList(
	c echo.Context,
	data templates.ListEditForm,
	listID string,
	adminID string,
) error {
	elements := make([]wishlister.WishListElement, len(data.Elements))

	for idx, elt := range data.Elements {
		elements[idx] = wishlister.WishListElement{
			Name:        elt.Name,
			Description: elt.Description,
			URL:         elt.URL,
		}
	}

	return s.wishlister.UpdateListElements(c.Request().Context(), listID, adminID, elements)
}

func listToEditData(list wishlister.WishList) templates.ListEditForm {
	data := templates.ListEditForm{Elements: make([]templates.ListEditFormElement, len(list.Elements))}
	for idx, element := range list.Elements {
		id, _ := nanoid.New()
		data.Elements[idx] = templates.ListEditFormElement{
			ID:          id,
			Name:        element.Name,
			Description: element.Description,
			URL:         element.URL,
		}
	}

	return data
}

func (s Server) validateElement(element templates.ListEditFormElement, ok bool) (templates.ListEditFormElement, bool) {
	if element.Name == "" {
		element.NameError = "Le nom ne peut pas être vide."
		ok = false
	} else if utf8.RuneCountInString(element.Name) > 255 {
		element.NameError = "Le nom ne peut pas dépasser 255 caractères."
		ok = false
	}

	if utf8.RuneCountInString(element.Description) > 500 {
		element.DescriptionError = "La description ne peut pas dépasser 500 caractères."
		ok = false
	}

	if element.URL != "" {
		if err := s.validate.Var(element.URL, "startswith=https://|startswith=http://,url"); err != nil {
			element.URLError = "L'URL n'est pas valide."
			ok = false
		} else if utf8.RuneCountInString(element.URL) > 2000 {
			element.URLError = "L'URL ne peut pas dépasser 2000 caractères."
			ok = false
		}
	}

	return element, ok
}
