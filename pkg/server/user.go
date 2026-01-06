package server

import (
	"errors"
	"net/http"

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

func (s Server) handleMagicLink(c echo.Context) error {
	sessionID := c.Param("token")
	session, err := s.wishlister.GetSession(c.Request().Context(), sessionID)
	if err != nil {
		s.e.Logger.Error("failed to get session from magic link: ", err)
		return renderOK(c, s.templates.RenderLogin, ParamsLogin{
			Error: "Le lien magique est invalide ou a expiré. Veuillez réessayer.",
		})
	}

	s.setUserSessionCookie(c, session.SessionID)
	return c.Redirect(302, "/lists")
}

func (s Server) setUserSessionCookie(c echo.Context, sessionID string) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}

	c.SetCookie(cookie)
}

func (s Server) getUserWishLists(c echo.Context) error {
	sessionIDCookie, err := c.Cookie("session_id")
	if err != nil {
		return c.Redirect(302, "/login")
	}

	session, err := s.wishlister.GetSession(c.Request().Context(), sessionIDCookie.Value)
	if err != nil {
		s.e.Logger.Error("failed to get session from cookie: ", err)
		return c.Redirect(302, "/login")
	}

	lists, err := s.wishlister.GetUserWishLists(c.Request().Context(), session.UserID)
	if err != nil {
		s.e.Logger.Error("failed to get user wish lists: ", err)
		return renderOK(c, s.templates.RenderUserListsView, ParamsUserListsView{
			Error: "Erreur lors de la récupération de vos listes de souhaits. Veuillez réessayer plus tard.",
		})
	}

	params := ParamsUserListsView{}
	for _, list := range lists {
		params.Lists = append(params.Lists, UserListsViewList{
			ID:      list.ID,
			AdminID: list.AdminID,
			Name:    list.Name,
		})
	}
	return renderOK(c, s.templates.RenderUserListsView, params)
}
