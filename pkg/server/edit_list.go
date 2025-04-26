package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"unicode/utf8"

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
		if errors.Is(err, wishlister.WishListNotFoundError{}) {
			return render(c, http.StatusNotFound, s.templates.RenderListNotFoundBytes, nil)
		}

		if errors.Is(err, wishlister.WishListInvalidAdminIDError{}) {
			return render(c, http.StatusForbidden, s.templates.RenderListAccessDeniedBytes, nil)
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
		c.Logger().Print("Error while binding wishlist elements form: %s", err)
		return data, false, ErrInvalidForm
	}

	// Validation at this step is very minimal: we only check that the 3 fields are
	// present. If this step is not OK, it probably means that the form wasn't sent
	// through the HTML page, and we just return an error.
	// The actual validation of the values of the fields is done later while building
	// the JSON that will be put on the page. This way, we can save error messages in
	// the JSON to show them on the HTML page.
	err = s.validate.Struct(form)
	if err != nil {
		c.Logger().Print("Error while validating wishlist elements form: %s", err)
		return data, false, ErrInvalidForm
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

		element, ok = s.validateElement(element, ok)
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

func (s Server) validateElement(element editListDataElement, ok bool) (editListDataElement, bool) {
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
