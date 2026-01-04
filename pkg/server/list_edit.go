package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/go-playground/form"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	nanoid "github.com/matoous/go-nanoid/v2"

	"github.com/erdnaxeli/wishlister"
)

type listEditTmplParams struct {
	Name string
	Data string
}

type editListForm struct {
	Elements []editListFormElement `json:"elements"`
}

type editListFormElement struct {
	ID string `json:"id"`

	Name             string `json:"name"`
	NameError        string `json:"name_error"`
	Description      string `json:"description"`
	DescriptionError string `json:"description_error"`
	URL              string `json:"url"`
	URLError         string `json:"url_error"`

	Error string `json:"error"`
}

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
		if errors.Is(err, wishlister.ErrWishListNotFound) {
			return render(c, http.StatusNotFound, s.templates.RenderListNotFound, nil)
		}

		if errors.Is(err, wishlister.ErrWishListInvalidAdminID) {
			return render(c, http.StatusForbidden, s.templates.RenderListAccessDenied, nil)
		}

		return err
	}

	var data editListForm

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

	return renderOK(c, s.templates.RenderListEdit, tmplParams)
}

func (s Server) validateEditForm(c echo.Context) (editListForm, bool, error) {
	data := editListForm{}
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
	form editListForm,
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

func listToEditData(list wishlister.WishList) editListForm {
	data := editListForm{Elements: make([]editListFormElement, len(list.Elements))}
	for idx, element := range list.Elements {
		id, _ := nanoid.New()
		data.Elements[idx] = editListFormElement{
			ID:          id,
			Name:        element.Name,
			Description: element.Description,
			URL:         element.URL,
		}
	}

	return data
}

func (s Server) validateElement(element editListFormElement, ok bool) (editListFormElement, bool) {
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
