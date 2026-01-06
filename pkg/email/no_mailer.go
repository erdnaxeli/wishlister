package email

import "log"

// NoMailer does not send any email.
type NoMailer struct{}

// SendNewWishListEmail actually does not send any email.
func (n NoMailer) SendNewWishListEmail(_ string, _ string, _ string, _ string) error {
	return nil
}

// SendMagicLink actually does not send any email.
func (n NoMailer) SendMagicLink(to string, sessionID string) error {
	log.Printf("NoMailer: SendMagicLink called for %s with sessionID %s", to, sessionID)
	return nil
}
