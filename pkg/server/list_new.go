package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/erdnaxeli/wishlister"
)

type createWishListForm struct {
	Name  string `form:"name"  validate:"required,max=255"`
	User  string `form:"user"  validate:"required,max=255"`
	Email string `form:"email" validate:"omitempty,email,max=255"`
}

func (s Server) getNewWishList(w http.ResponseWriter, r *http.Request) {
	params := ParamsNew{}

	sessionIDCookie, err := r.Cookie("session_id")
	if err == nil {
		session, err := s.wishlister.GetSession(r.Context(), sessionIDCookie.Value)
		if err == nil {
			params.User = session.Username
			params.Email = session.UserEmail
		}
	}

	s.renderOK(w, s.templates.RenderNew, params)
}

func (s Server) createNewWishList(w http.ResponseWriter, r *http.Request) {
	form := createWishListForm{
		Name:  r.PostFormValue("name"),
		User:  r.PostFormValue("user"),
		Email: r.PostFormValue("email"),
	}

	err := s.validate.Struct(form)
	if err != nil {
		s.handleNewWishListError(w, form, err)
		return
	}

	params := wishlister.CreateWishlistParams{
		Name:      form.Name,
		Username:  form.User,
		UserEmail: form.Email,
	}

	listID, adminID, err := s.wishlister.CreateWishList(
		r.Context(),
		params,
	)
	if err != nil {
		// We could get errors about empty fields here, but this should have been catched
		// by the validation step before. So we just log and return a generic error.
		s.logger.Error("failed to create new wish list: ", "err", err)
		s.renderOK(w, s.templates.RenderNew, ParamsNew{
			Error: "Erreur lors de la soumission du formulaire, veuillez réessayer",
			Name:  form.Name,
			User:  form.User,
			Email: form.Email,
		})
		return
	}

	http.Redirect(w, r, fmt.Sprintf("l/%s/%s", listID, adminID), http.StatusSeeOther)
}

func (s Server) handleNewWishListError(
	w http.ResponseWriter,
	form createWishListForm,
	err error,
) {
	formError := ParamsNew{
		Name:  form.Name,
		User:  form.User,
		Email: form.Email,
	}

	var invalidErr *validator.InvalidValidationError
	if errors.As(err, &invalidErr) {
		s.logger.Error("invalid validation error", "err", err)

		formError.Error = "Erreur lors de la soumission du formulaire, veuillez réessayer."
		s.renderOK(w, s.templates.RenderNew, formError)
		return
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		s.handleNewWishListValidationsErrors(w, formError, validationErrors)
		return
	}

	s.logger.Error("unknown error during form validation: ", "err", err)
	formError.Error = "Erreur lors de la soumission du formulaire, veuillez réessayer."
	s.renderOK(w, s.templates.RenderNew, formError)
}

func (s Server) handleNewWishListValidationsErrors(
	w http.ResponseWriter,
	formError ParamsNew,
	validationErrors validator.ValidationErrors,
) {
	if len(validationErrors) == 0 {
		// that should not happend
		s.logger.Error("validation errors but length is 0")
		formError.Error = "Erreur lors de la soumission du formulaire, veuillez réessayer."
		s.renderOK(w, s.templates.RenderNew, formError)
		return
	}

	validationErr := validationErrors[0]
	switch validationErr.Field() {
	case "Name":
		s.handleNewWishListNameError(w, formError, validationErr)
		return
	case "User":
		s.handleNewWishListUserError(w, formError, validationErr)
		return
	case "Email":
		s.handleNewWishListEmailError(w, formError, validationErr)
		return
	default:
		s.logger.Error("unknown validation error field", "field", validationErr.Field())
		s.renderOK(w, s.templates.RenderNew, formError)
		return
	}
}

func (s Server) handleNewWishListNameError(
	w http.ResponseWriter,
	formError ParamsNew,
	validationErr validator.FieldError,
) {
	switch validationErr.Tag() {
	case "required":
		formError.NameError = "Le nom est requis."
	case "max":
		formError.NameError = "Le nom doit faire moins de 255 caractères."
	default:
		s.logger.Error(
			"unknown validation error tag on name field",
			"tag", validationErr.Tag(),
		)
		formError.NameError = "Le nom est requis et doit faire moins de 255 caractères."
	}

	s.renderOK(w, s.templates.RenderNew, formError)
}

func (s Server) handleNewWishListUserError(
	w http.ResponseWriter,
	formError ParamsNew,
	validationErr validator.FieldError,
) {
	switch validationErr.Tag() {
	case "required":
		formError.UserError = "Le nom d'utilisateur est requis."
	case "max":
		formError.UserError = "Le nom d'utilisateur doit faire moins de 255 caractères."
	default:
		s.logger.Error(
			"unknown validation error tag on user field",
			"tag", validationErr.Tag(),
		)
		formError.UserError = "Le nom d'utilisateur est requis et doit faire moins de 255 caractères."
	}

	s.renderOK(w, s.templates.RenderNew, formError)
}

func (s Server) handleNewWishListEmailError(
	w http.ResponseWriter,
	formError ParamsNew,
	validationErr validator.FieldError,
) {
	switch validationErr.Tag() {
	case "email":
		formError.EmailError = "L'adresse email n'est pas valide."
	case "max":
		formError.EmailError = "L'adresse email doit faire moins de 255 caractères."
	default:
		s.logger.Error(
			"unknown validation error tag on email field: ",
			"tag", validationErr.Tag(),
		)
		formError.EmailError = "L'adresse email est requise et doit faire moins de 255 caractères."
	}

	s.renderOK(w, s.templates.RenderNew, formError)
}
