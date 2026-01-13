package email

import (
	"context"
	"fmt"

	"github.com/go-hermes/hermes/v2"
	"github.com/wneessen/go-mail"
)

func (s smtpSender) SendNewWishListEmail(
	ctx context.Context,
	to string,
	username string,
	listID string,
	adminID string,
) error {
	htmlBody, textBody, err := s.getNewWishListMailBody(username, listID, adminID)
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

	mailMsg.Subject("Liste de voeux créé")
	mailMsg.SetBodyString(mail.TypeTextHTML, htmlBody)
	mailMsg.AddAlternativeString(mail.TypeTextPlain, textBody)

	err = s.client.DialAndSendWithContext(ctx, mailMsg)
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
