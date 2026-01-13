package email

import (
	"context"
	"fmt"

	"github.com/go-hermes/hermes/v2"
	"github.com/wneessen/go-mail"
)

func (s smtpSender) SendMagicLink(
	ctx context.Context,
	to string,
	magicLink string,
) error {
	htmlBody, textBody, err := s.getMagicLinkMailBody(magicLink)
	if err != nil {
		return err
	}

	mailMsg := mail.NewMsg()
	err = mailMsg.From(s.from)
	if err != nil {
		return err
	}
	err = mailMsg.To(to)
	if err != nil {
		return err
	}
	mailMsg.Subject("Votre lien de connexion")
	mailMsg.SetBodyString(mail.TypeTextHTML, htmlBody)
	mailMsg.AddAlternativeString(mail.TypeTextPlain, textBody)

	err = s.client.DialAndSendWithContext(ctx, mailMsg)
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
							"https://www.malistedevoeux.fr/login/magic/%s",
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
