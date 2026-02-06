package server

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type sendMagicLinkForm struct {
	Email string `form:"email" validate:"required,email,max=255"`
}

func (s Server) sendMagicLink(w http.ResponseWriter, r *http.Request) {
	form := sendMagicLinkForm{
		Email: r.PostFormValue("email"),
	}

	err := s.validate.Struct(form)
	if err != nil {
		s.handleSendMagicLinkError(w, form, err)
		return
	}

	err = s.wishlister.SendMagicLink(r.Context(), form.Email)
	if err != nil {
		s.logger.Error("failed to send magic link", "err", err)

		s.renderOK(w, s.templates.RenderLogin, ParamsLogin{
			Error: "Erreur lors de l'envoi du lien magique, veuillez réessayer.",
			Email: form.Email,
		})
		return
	}

	s.renderOK(w, s.templates.RenderLogin, ParamsLogin{
		Email: form.Email,
		Sent:  true,
	})
}

func (s Server) handleSendMagicLinkError(
	w http.ResponseWriter,
	form sendMagicLinkForm,
	validationErr error,
) {
	params := ParamsLogin{
		Email: form.Email,
	}

	var invalidErr *validator.InvalidValidationError
	if errors.As(validationErr, &invalidErr) {
		s.logger.Error("invalid validation error", "err", invalidErr)

		params.Error = "Erreur lors de la soumission du formulaire, veuillez réessayer."
		s.renderOK(w, s.templates.RenderLogin, params)
	}

	var validationErrors validator.ValidationErrors
	if errors.As(validationErr, &validationErrors) {
		if len(validationErrors) == 0 {
			// that should not happen
			s.logger.Error("validation errors but length is 0")
			params.Error = "Erreur lors de la soumission du formulaire, veuillez réessayer."
			s.renderOK(w, s.templates.RenderLogin, params)
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
				s.logger.Error("unknown email validation error tag", "tag", validationErr.Tag())
			}
		default:
			s.logger.Error("unknown validation error field", "field", validationErr.Field())
		}

		s.renderOK(w, s.templates.RenderLogin, params)
	}

	s.logger.Error("unknown error during form validation", "err", validationErr)
	params.Error = "Erreur lors de la soumission du formulaire, veuillez réessayer."
	s.renderOK(w, s.templates.RenderLogin, params)
}

func (s Server) handleMagicLink(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	session, err := s.wishlister.GetSessionByMagicLink(r.Context(), token)
	if err != nil {
		s.logger.Error("failed to get session from magic link", "err", err)
		s.renderOK(w, s.templates.RenderLogin, ParamsLogin{
			Error: "Le lien magique est invalide ou a expiré. Veuillez réessayer.",
		})
		return
	}

	s.setUserSessionCookie(w, session.SessionID)
	http.Redirect(w, r, "/lists", http.StatusFound)
}

func (s Server) setUserSessionCookie(w http.ResponseWriter, sessionID string, expire ...bool) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}

	if len(expire) > 0 && expire[0] {
		cookie.MaxAge = -1
	}

	http.SetCookie(w, cookie)
}

func (s Server) getUserWishLists(w http.ResponseWriter, r *http.Request) {
	sessionIDCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	session, err := s.wishlister.GetSession(r.Context(), sessionIDCookie.Value)
	if err != nil {
		s.logger.Error("failed to get session from cookie", "err", err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	lists, err := s.wishlister.GetUserWishLists(r.Context(), session.UserID)
	if err != nil {
		s.logger.Error("failed to get user wish lists: ", "err", err)
		s.renderOK(w, s.templates.RenderUserListsView, ParamsUserListsView{
			Error: "Erreur lors de la récupération de vos listes de souhaits. Veuillez réessayer plus tard.",
		})
		return
	}

	params := ParamsUserListsView{}
	for _, list := range lists {
		params.Lists = append(params.Lists, UserListsViewList{
			ID:      list.ID,
			AdminID: list.AdminID,
			Name:    list.Name,
		})
	}

	s.renderOK(w, s.templates.RenderUserListsView, params)
}

func (s Server) logout(w http.ResponseWriter, r *http.Request) {
	sessionIDCookie, err := r.Cookie("session_id")
	if err == nil {
		s.wishlister.DeleteSession(r.Context(), sessionIDCookie.Value)
	}

	s.setUserSessionCookie(w, "", true)
	s.renderOK(w, s.templates.RenderLogout, nil)
}
