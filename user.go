package wishlister

import (
	"context"
	"database/sql"
	"errors"
	"log"

	nanoid "github.com/matoous/go-nanoid/v2"

	"github.com/erdnaxeli/wishlister/pkg/repository"
)

// GetOrCreateUser retrieves an existing user by email or creates a new one.
//
// It returns the user ID.
func (a *app) GetOrCreateUser(
	ctx context.Context,
	username string,
	email string,
) (string, error) {
	userID, _ := nanoid.New()

	user, err := a.queries.GetOrCreateUser(ctx, repository.GetOrCreateUserParams{
		ID:    userID,
		Name:  username,
		Email: email,
	})
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (a *app) createUserSession(
	ctx context.Context,
	userID string,
) (Session, error) {
	sessionID, _ := nanoid.New()
	magicLinkToken, _ := nanoid.New()
	err := a.queries.CreateUserSession(ctx, repository.CreateUserSessionParams{
		ID:             sessionID,
		UserID:         userID,
		MagicLinkToken: NewNullString(magicLinkToken),
	})
	if err != nil {
		return Session{}, err
	}

	return Session{
		UserID:         userID,
		SessionID:      sessionID,
		MagicLinkToken: magicLinkToken,
	}, nil
}

func (a *app) SendMagicLink(ctx context.Context, email string) error {
	userID, err := a.GetOrCreateUser(ctx, "", email)
	if err != nil {
		return err
	}

	session, err := a.createUserSession(ctx, userID)
	if err != nil {
		return err
	}

	return a.emailSender.SendMagicLink(ctx, email, session.MagicLinkToken)
}

func (a *app) GetSession(ctx context.Context, sessionID string) (Session, error) {
	session, err := a.queries.GetUserSession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Session{}, ErrSessionNotFound
		}

		return Session{}, err
	}

	return Session{
		UserID:    session.UserID,
		Username:  session.Username,
		UserEmail: session.UserEmail,
		SessionID: session.ID,
	}, nil
}

// GetSessionByMagicLink returns the session associated with the given magic link token.
//
// If the session is not found, an error ErrSessionNotFound is returned.
// Once used, the magic link token is invalidated and cannot be used again.
func (a *app) GetSessionByMagicLink(ctx context.Context, token string) (Session, error) {
	session, err := a.queries.GetUserSessionByMagicLink(ctx, NewNullString(token))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Session{}, ErrSessionNotFound
		}

		return Session{}, err
	}

	return Session{
		UserID:    session.UserID,
		SessionID: session.ID,
	}, nil
}

func (a *app) DeleteSession(ctx context.Context, sessionID string) {
	err := a.queries.DeleteUserSession(ctx, sessionID)
	if err != nil {
		// We just log the error, as the user can still be logged out even if the session
		// still exists.
		log.Printf("Error while deleting session %s", sessionID)
	}
}
