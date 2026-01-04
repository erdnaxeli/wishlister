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

// GetOrCreateUserSession retrieves an existing session by sessionID or creates a new one.
//
// If sessionID is empty, a new session is created.
// If sessionID is not empty and not found, a new session is created.
// If sessionID is not empty and found:
//   - if the session belongs to the user, nothing is done.
//   - if the session does not belong to the user, a new session is created and the old one is deleted.
//
// It returns the (same or newly created) session ID.
func (a *app) GetOrCreateUserSession(
	ctx context.Context,
	userID string,
	sessionID string,
) (string, error) {
	if sessionID != "" {
		return a.createUserSession(ctx, userID)
	}

	session, err := a.queries.GetUserSession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return a.createUserSession(ctx, userID)
		}

		return "", err
	}

	if session.UserID == userID {
		return sessionID, nil
	}

	err = a.queries.DeleteUserSession(ctx, sessionID)
	if err != nil {
		// Log the error but continue as we can create a new session anyway
		log.Printf("error while deleting user session %s", sessionID)
	}

	return a.createUserSession(ctx, userID)
}

func (a *app) createUserSession(
	ctx context.Context,
	userID string,
) (string, error) {
	sessionID, _ := nanoid.New()
	err := a.queries.CreateUserSession(ctx, repository.CreateUserSessionParams{
		ID:     sessionID,
		UserID: userID,
	})
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (a *app) SendMagicLink(ctx context.Context, email string) error {
	userID, err := a.GetOrCreateUser(ctx, "", email)
	if err != nil {
		return err
	}

	sessionID, err := a.createUserSession(ctx, userID)
	if err != nil {
		return err
	}

	return a.emailSender.SendMagicLink(email, sessionID)
}
