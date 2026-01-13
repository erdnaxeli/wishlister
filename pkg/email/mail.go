// Package email implement methods to send emails.
package email

import (
	"context"
	"time"

	"github.com/go-hermes/hermes/v2"
	"github.com/wneessen/go-mail"
)

// Sender is the main interface of this package.
type Sender interface {
	// SendNewWishListEmail send a mail for a newly created wishlist.
	//
	// The mail contains the link to share, and the link to edit the wishlist.
	SendNewWishListEmail(
		ctx context.Context,
		to string,
		username string,
		listID string,
		adminID string,
	) error

	// SendMagicLink sends a magic link to the given email address.
	//
	// The link can be used to login the user.
	SendMagicLink(ctx context.Context, to string, sessionID string) error
}

type smtpSender struct {
	client *mail.Client
	from   string

	h hermes.Hermes
}

// NewSMTPSender return a Sender object.
func NewSMTPSender(
	username string,
	password string,
	server string,
	port int,
	from string,
) (Sender, error) {
	client, err := mail.NewClient(
		server,
		mail.WithPort(port),
		mail.WithSSL(),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithTimeout(10*time.Second),
	)
	if err != nil {
		return smtpSender{}, err
	}

	return smtpSender{
		client: client,
		from:   from,

		h: hermes.Hermes{
			Product: hermes.Product{
				Name:        "Ma liste de vœux .fr",
				Link:        "https://www.malistdevoeux.fr/",
				Copyright:   " ",
				TroubleText: "Si le bouton {ACTION} ne marche pas, copier l'URL suivante dans votre navigateur.",
			},
		},
	}, nil
}
