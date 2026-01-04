package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/erdnaxeli/wishlister"
)

type createWishListForm struct {
	Name  string `form:"name"  validate:"required,max=255"`
	User  string `form:"user"  validate:"required,max=255"`
	Email string `form:"email" validate:"required,email,max=255"`
}

func (s Server) createNewWishList(c echo.Context) error {
	form := createWishListForm{}
	err := c.Bind(&form)
	if err != nil {
		s.e.Logger.Error("failed to bind create wish list form: ", err)

		return renderOK(c, s.templates.RenderNew, ParamsNew{
			Error: "Erreur lors de la soumission du formulaire, veuillez réessayer.",
			Name:  form.Name,
			User:  form.User,
			Email: form.Email,
		})
	}

	err = s.validate.Struct(form)
	if err != nil {
		return s.handleNewWishListError(c, form, err)
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
		// We could get errors about empty fields here, but this should have been catched
		// by the validation step before. So we just log and return a generic error.
		s.e.Logger.Error("failed to create new wish list: ", err)
		return renderOK(c, s.templates.RenderNew, ParamsNew{
			Error: "Erreur lors de la soumission du formulaire, veuillez réessayer",
			Name:  form.Name,
			User:  form.User,
			Email: form.Email,
		})
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("l/%s/%s", listID, adminID))
}

func (s Server) handleNewWishListError(
	c echo.Context,
	form createWishListForm,
	err error,
) error {
	formError := ParamsNew{
		Name:  form.Name,
		User:  form.User,
		Email: form.Email,
	}

	var invalidErr *validator.InvalidValidationError
	if errors.As(err, &invalidErr) {
		s.e.Logger.Error("invalid validation error: ", err)

		formError.Error = "Erreur lors de la soumission du formulaire, veuillez réessayer."
		return renderOK(c, s.templates.RenderNew, formError)
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return s.handleNewWishListValidationsErrors(c, formError, validationErrors)
	}

	s.e.Logger.Error("unknown error during form validation: ", err)
	formError.Error = "Erreur lors de la soumission du formulaire, veuillez réessayer."
	return renderOK(c, s.templates.RenderNew, formError)
}

func (s Server) handleNewWishListValidationsErrors(
	c echo.Context,
	formError ParamsNew,
	validationErrors validator.ValidationErrors,
) error {
	if len(validationErrors) == 0 {
		// that should not happend
		s.e.Logger.Error("validation errors but length is 0")
		formError.Error = "Erreur lors de la soumission du formulaire, veuillez réessayer."
		return renderOK(c, s.templates.RenderNew, formError)
	}

	validationErr := validationErrors[0]
	switch validationErr.Field() {
	case "Name":
		return s.handleNewWishListNameError(c, formError, validationErr)
	case "User":
		return s.handleNewWishListUserError(c, formError, validationErr)
	case "Email":
		return s.handleNewWishListEmailError(c, formError, validationErr)
	default:
		s.e.Logger.Error("unknown validation error field: ", validationErr.Field())
		return renderOK(c, s.templates.RenderNew, formError)
	}
}

func (s Server) handleNewWishListNameError(
	c echo.Context,
	formError ParamsNew,
	validationErr validator.FieldError,
) error {
	switch validationErr.Tag() {
	case "required":
		formError.NameError = "Le nom est requis."
	case "max":
		formError.NameError = "Le nom doit faire moins de 255 caractères."
	default:
		s.e.Logger.Error(
			"unknown validation error tag on name field: ",
			validationErr.Tag(),
		)
		formError.NameError = "Le nom est requis et doit faire moins de 255 caractères."
	}

	return renderOK(c, s.templates.RenderNew, formError)
}

func (s Server) handleNewWishListUserError(
	c echo.Context,
	formError ParamsNew,
	validationErr validator.FieldError,
) error {
	switch validationErr.Tag() {
	case "required":
		formError.UserError = "Le nom d'utilisateur est requis."
	case "max":
		formError.UserError = "Le nom d'utilisateur doit faire moins de 255 caractères."
	default:
		s.e.Logger.Error(
			"unknown validation error tag on user field: ",
			validationErr.Tag(),
		)
		formError.UserError = "Le nom d'utilisateur est requis et doit faire moins de 255 caractères."
	}

	return renderOK(c, s.templates.RenderNew, formError)
}

func (s Server) handleNewWishListEmailError(
	c echo.Context,
	formError ParamsNew,
	validationErr validator.FieldError,
) error {
	switch validationErr.Tag() {
	case "required":
		formError.EmailError = "L'adresse email est requise."
	case "email":
		formError.EmailError = "L'adresse email n'est pas valide."
	case "max":
		formError.EmailError = "L'adresse email doit faire moins de 255 caractères."
	default:
		s.e.Logger.Error(
			"unknown validation error tag on email field: ",
			validationErr.Tag(),
		)
		formError.EmailError = "L'adresse email est requise et doit faire moins de 255 caractères."
	}

	return renderOK(c, s.templates.RenderNew, formError)
}
