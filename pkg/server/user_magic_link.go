package server

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type sendMagicLinkForm struct {
	Email string `form:"email" validate:"required,email,max=255"`
}

func (s Server) sendMagicLink(c echo.Context) error {
	form := sendMagicLinkForm{}
	err := c.Bind(&form)
	if err != nil {
		s.e.Logger.Error("failed to bind send magic link form: ", err)

		return renderOK(c, s.templates.RenderLogin, ParamsLogin{
			Error: "Erreur lors de la soumission du formulaire, veuillez réessayer.",
			Email: form.Email,
		})
	}

	err = s.validate.Struct(form)
	if err != nil {
		return s.handleSendMagicLinkError(c, form, err)
	}

	err = s.wishlister.SendMagicLink(c.Request().Context(), form.Email)
	if err != nil {
		s.e.Logger.Error("failed to send magic link: ", err)

		return renderOK(c, s.templates.RenderLogin, ParamsLogin{
			Error: "Erreur lors de l'envoi du lien magique, veuillez réessayer.",
			Email: form.Email,
		})
	}

	return renderOK(c, s.templates.RenderLogin, ParamsLogin{
		Email: form.Email,
		Sent:  true,
	})
}

func (s Server) handleSendMagicLinkError(
	c echo.Context,
	form sendMagicLinkForm,
	validationErr error,
) error {
	params := ParamsLogin{
		Email: form.Email,
	}

	var invalidErr *validator.InvalidValidationError
	if errors.As(validationErr, &invalidErr) {
		s.e.Logger.Error("invalid validation error: ", invalidErr)

		params.Error = "Erreur lors de la soumission du formulaire, veuillez réessayer."
		return renderOK(c, s.templates.RenderLogin, params)
	}

	var validationErrors validator.ValidationErrors
	if errors.As(validationErr, &validationErrors) {
		if len(validationErrors) == 0 {
			// that should not happen
			s.e.Logger.Error("validation errors but length is 0")
			params.Error = "Erreur lors de la soumission du formulaire, veuillez réessayer."
			return renderOK(c, s.templates.RenderLogin, params)
		}

		validationErr := validationErrors[0]
		switch validationErr.Field() {
		case "Email":
			switch validationErr.Tag() {
			case "required":
				params.EmailError = "L'adresse email est requise."
			case "email":
				params.EmailError = "L'adresse email n'est pas valide."
			case "max":
				params.EmailError = "L'adresse email est trop longue."
			default:
				s.e.Logger.Error("unknown email validation error tag: ", validationErr.Tag())
			}
		default:
			s.e.Logger.Error("unknown validation error field: ", validationErr.Field())
		}

		return renderOK(c, s.templates.RenderLogin, params)
	}

	s.e.Logger.Error("unknown error during form validation: ", validationErr)
	params.Error = "Erreur lors de la soumission du formulaire, veuillez réessayer."
	return renderOK(c, s.templates.RenderLogin, params)
}
