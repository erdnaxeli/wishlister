package server

import (
	"fmt"
	"net/http"

	"github.com/erdnaxeli/wishlister"
)

type createGroupForm struct {
	Name  string `form:"name"  validate:"required,max=255"`
	Email string `form:"email" validate:"required,email,max=255"`
}

func (s Server) createNewGroup(w http.ResponseWriter, r *http.Request) {
	form := createGroupForm{
		Name:  r.PostFormValue("name"),
		Email: r.PostFormValue("email"),
	}

	err := s.validate.Struct(form)
	if err != nil {
		s.logger.Error("validation error", "err", err)
		return
	}

	groupID, err := s.wishlister.CreateGroup(
		r.Context(),
		wishlister.CreateGroupParams{
			Name:      form.Name,
			UserEmail: form.Email,
		},
	)
	if err != nil {
		s.logger.Error("error while creating group", "err", err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/g/%s", groupID), http.StatusSeeOther)
}
