package email

import (
	"fmt"

	"github.com/go-hermes/hermes/v2"
	"gopkg.in/gomail.v2"
)

func (s smtpSender) SendMagicLink(
	to string,
	magicLink string,
) error {
	htmlBody, textBody, err := s.getMagicLinkMailBody(magicLink)
	if err != nil {
		return err
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", s.from)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", "Votre lien de connexion")
	mail.SetBody("text/html", htmlBody)
	mail.AddAlternative("text/plain", textBody)

	err = s.dialer.DialAndSend(mail)
	if err != nil {
		return err
	}

	return nil
}

func (s smtpSender) getMagicLinkMailBody(
	sessionID string,
) (string, string, error) {
	mail := hermes.Email{
		Body: hermes.Body{
			Greeting: "Bonjour",
			Intros: []string{
				"Voici votre lien de connexion :",
			},
			Actions: []hermes.Action{
				{
					Button: hermes.Button{
						Text: "Se connecter",
						Link: fmt.Sprintf(
							"https://www.malistedevoeux.fr/login/%s",
							sessionID,
						),
					},
				},
			},
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
