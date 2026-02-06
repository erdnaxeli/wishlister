package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/erdnaxeli/wishlister"
)

type getWishListParam struct {
	AdminID string `param:"adminID"`
	ListID  string `param:"listID"`
}

func readWishListParam(r *http.Request) getWishListParam {
	return getWishListParam{
		AdminID: chi.URLParam(r, "adminID"),
		ListID:  chi.URLParam(r, "listID"),
	}
}

func (s Server) getWishList(w http.ResponseWriter, r *http.Request) {
	params := readWishListParam(r)

	var list wishlister.WishList
	var err error

	if params.AdminID == "" {
		list, err = s.wishlister.GetWishList(r.Context(), params.ListID)
		if err != nil {
			s.logger.Error("error while getting wishlist", "err", err)
			panic(err)
		}
	} else {
		list, err = s.wishlister.GetEditableWishList(r.Context(), params.ListID, params.AdminID)
		if err != nil {
			if errors.Is(err, wishlister.ErrWishListInvalidAdminID) {
				http.Redirect(w, r, fmt.Sprintf("/%s", params.ListID), http.StatusMovedPermanently)
				return
			}

			panic(err)
		}
	}

	s.renderOK(w, s.templates.RenderListView, list)
}
