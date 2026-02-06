package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/go-playground/form"
	"github.com/google/uuid"
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

func (s Server) editList(w http.ResponseWriter, r *http.Request) {
	params := readWishListParam(r)

	list, err := s.wishlister.GetEditableWishList(
		r.Context(),
		params.ListID,
		params.AdminID,
	)
	if err != nil {
		if errors.Is(err, wishlister.ErrWishListNotFound) {
			s.render(w, http.StatusNotFound, s.templates.RenderListNotFound, nil)
			return
		}

		if errors.Is(err, wishlister.ErrWishListInvalidAdminID) {
			s.render(w, http.StatusForbidden, s.templates.RenderListAccessDenied, nil)
			return
		}

		panic(err)
	}

	var data editListForm

	if r.Method == http.MethodPost {
		var ok bool
		data, ok, err = s.validateEditForm(r)
		if err != nil {
			panic(err)
		}

		if ok {
			err := s.updateList(r, data, params.ListID, params.AdminID)
			if err != nil {
				panic(err)
			}

			http.Redirect(
				w, r,
				fmt.Sprintf("/l/%s/%s", params.ListID, params.AdminID),
				http.StatusSeeOther,
			)
			return
		}
	} else {
		data = listToEditData(list)
	}

	dataJSON, err := json.Marshal(data.Elements)
	if err != nil {
		panic(err)
	}

	tmplParams := listEditTmplParams{
		Name: list.Name,
		Data: string(dataJSON),
	}

	s.renderOK(w, s.templates.RenderListEdit, tmplParams)
}

func (s Server) validateEditForm(
	r *http.Request,
) (editListForm, bool, error) {
	data := editListForm{}
	ok := true
	decoder := form.NewDecoder()
	err := r.ParseForm()
	if err != nil {
		s.logger.Error("Error while reading wishlist elements form", "err", err)
		return data, false, err
	}

	err = decoder.Decode(&data, r.Form)
	if err != nil {
		s.logger.Error("Error while decoding wishlist elements form", "err", err)
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
	r *http.Request,
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

	return s.wishlister.UpdateListElements(r.Context(), listID, adminID, elements)
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
		err := s.validate.Var(element.URL, "startswith=https://|startswith=http://,url")
		if err != nil {
			element.URLError = "L'URL n'est pas valide."
			ok = false
		} else if utf8.RuneCountInString(element.URL) > 2000 {
			element.URLError = "L'URL ne peut pas dépasser 2000 caractères."
			ok = false
		}
	}

	return element, ok
}
