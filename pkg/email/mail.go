// Package email implement methods to send emails.
package email

import (
	"github.com/go-hermes/hermes/v2"
	"gopkg.in/gomail.v2"
)

// Sender is the main interface of this package.
type Sender interface {
	// SendNewWishListEmail send a mail for a newly created wishlist.
	//
	// The mail contains the link to share, and the link to edit the wishlist.
	SendNewWishListEmail(to string, username string, listID string, adminID string) error

	// SendMagicLink sends a magic link to the given email address.
	//
	// The link can be used to login the user.
	SendMagicLink(to string, sessionID string) error
}

type smtpSender struct {
	dialer *gomail.Dialer
	from   string

	h hermes.Hermes
}

// NewSMTPSender return a Sender object.
func NewSMTPSender(username string, password string, server string, port int, from string) Sender {
	return smtpSender{
		dialer: gomail.NewDialer(server, port, username, password),
		from:   from,

		h: hermes.Hermes{
			Product: hermes.Product{
				Name:        "Ma liste de vœux .fr",
				Link:        "https://www.malistdevoeux.fr/",
				Copyright:   " ",
				TroubleText: "Si le bouton {ACTION} ne marche pas, copier l'URL suivante dans votre navigateur.",
			},
		},
	}
}
