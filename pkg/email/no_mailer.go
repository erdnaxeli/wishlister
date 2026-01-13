package email

import (
	"context"
	"log"
)

// NoMailer does not send any email.
type NoMailer struct{}

// SendNewWishListEmail actually does not send any email.
func (n NoMailer) SendNewWishListEmail(
	_ context.Context,
	_ string,
	_ string,
	_ string,
	_ string,
) error {
	return nil
}

// SendMagicLink actually does not send any email.
func (n NoMailer) SendMagicLink(_ context.Context, to string, sessionID string) error {
	log.Printf("NoMailer: SendMagicLink called for %s with sessionID %s", to, sessionID)
	return nil
}
