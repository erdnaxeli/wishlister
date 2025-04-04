package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/erdnaxeli/wishlister"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	nanoid "github.com/matoous/go-nanoid/v2"
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

func editList(c echo.Context, app wishlister.App) error {
	listID := c.Param("listID")
	adminID := c.Param("adminID")

	list, err := app.GetEditableWishList(c.Request().Context(), listID, adminID)
	if err != nil {
		if errors.Is(err, wishlister.WishListNotFoundError{}) {
			return c.Render(http.StatusNotFound, "listNotFound", nil)
		}

		if errors.Is(err, wishlister.WishListInvalidAdminIDError{}) {
			return c.Render(http.StatusForbidden, "listAccessDenied", list)
		}

		return err
	}

	var dataJSON []byte

	if c.Request().Method == http.MethodPost {
		form, ok, err := validateEditForm(c)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		if ok {
			err := updateList(c, app, form, listID, adminID)
			if err != nil {
				return err
			}

			return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/%s/%s", listID, adminID))
		}

		dataJSON, err = json.Marshal(form.Elements)
		if err != nil {
			return err
		}
	} else {
		form := listToEditForm(list)
		dataJSON, err = json.Marshal(form.Elements)
		if err != nil {
			return err
		}
	}

	params := listEditTmplParams{
		Name: list.Name,
		Data: string(dataJSON),
	}

	return c.Render(http.StatusOK, "listEdit", params)
}

func validateEditForm(c echo.Context) (editListForm, bool, error) {
	form := editListForm{}
	ok := true

	values, _ := c.FormParams()
	nameValues, nameOk := values["name"]
	descriptionValues, descriptionOk := values["description"]
	urlValues, urlOk := values["url"]

	if !nameOk || !descriptionOk || !urlOk ||
		len(nameValues) != len(descriptionValues) || len(nameValues) != len(urlValues) {
		return form, false, ErrInvalidForm
	}

	for i := range nameValues {
		element := editListFormElement{
			ID:          uuid.NewString(),
			Name:        nameValues[i],
			Description: descriptionValues[i],
			URL:         urlValues[i],
		}

		if nameValues[i] == "" {
			element.NameError = "Le nom ne peut pas être vide."
			ok = false
		}

		form.Elements = append(form.Elements, element)
	}

	return form, ok, nil
}

func updateList(c echo.Context, app wishlister.App, form editListForm, listID string, adminID string) error {
	elements := make([]wishlister.WishListElement, len(form.Elements))

	for idx, elt := range form.Elements {
		elements[idx] = wishlister.WishListElement{
			Name:        elt.Name,
			Description: elt.Description,
			URL:         elt.URL,
		}
	}

	return app.UpdateListElements(c.Request().Context(), listID, adminID, elements)
}

func listToEditForm(list wishlister.WishList) editListForm {
	form := editListForm{Elements: make([]editListFormElement, len(list.Elements))}
	for idx, element := range list.Elements {
		id, _ := nanoid.New()
		form.Elements[idx] = editListFormElement{
			ID:          id,
			Name:        element.Name,
			Description: element.Description,
			URL:         element.URL,
		}
	}

	return form
}
