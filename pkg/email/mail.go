// Package email implement methods to send emails.
package email

import (
	"fmt"

	"github.com/go-hermes/hermes/v2"
	"gopkg.in/gomail.v2"
)

// Sender is the main interface of this package.
type Sender interface {
	// SendNewWishListEmail send a mail for a newly created wishlist.
	//
	// The mail contains the link to share, and the link to edit the wishlist.
	SendNewWishListEmail(to string, username string, listID string, adminID string) error
}

// NoMailer does not send any email.
type NoMailer struct{}

// SendNewWishListEmail actually does not send any email.
func (n NoMailer) SendNewWishListEmail(_ string, _ string, _ string, _ string) error {
	return nil
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

func (s smtpSender) SendNewWishListEmail(
	to string,
	username string,
	listID string,
	adminID string,
) error {
	htmlBody, textBody, err := s.getNewWishListMailBody(username, listID, adminID)
	if err != nil {
		return err
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", s.from)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", "Liste de voeux créé")
	mail.SetBody("text/html", htmlBody)
	mail.AddAlternative("text/plain", textBody)

	err = s.dialer.DialAndSend(mail)
	if err != nil {
		return err
	}

	return nil
}

func (s smtpSender) getNewWishListMailBody(
	username string,
	listID string,
	adminID string,
) (string, string, error) {
	mail := hermes.Email{
		Body: hermes.Body{
			Name:     username,
			Greeting: "Bonjour",
			Intros: []string{
				"Votre liste de vœux a bien été créé !",
				"Voici le lien à partager :",
				fmt.Sprintf("https://www.malistedevoeux.fr/%s", listID),
			},
			Actions: []hermes.Action{
				{
					Button: hermes.Button{
						Text: "Éditer la liste",
						Link: fmt.Sprintf(
							"https://www.malistedevoeux.fr/%s/%s/edit",
							listID,
							adminID,
						),
					},
				},
			},
			Signature: "À bientôt",
		},
	}

	htmlBody, err := s.h.GenerateHTML(mail)
	if err != nil {
		return "", "", err
	}

	textBody, err := s.h.GeneratePlainText(mail)
	if err != nil {
		return "", "", err
	}

	return htmlBody, textBody, nil
}
