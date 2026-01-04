package email

import (
	"fmt"

	"github.com/go-hermes/hermes/v2"
	"gopkg.in/gomail.v2"
)

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
