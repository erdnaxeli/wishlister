package wishlister

import (
	"context"

	nanoid "github.com/matoous/go-nanoid/v2"

	"github.com/erdnaxeli/wishlister/pkg/repository"
)

func (a *app) getOrCreateUser(
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
