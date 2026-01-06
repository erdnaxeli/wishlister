package wishlister

import (
	"context"
	"log"

	nanoid "github.com/matoous/go-nanoid/v2"

	"github.com/erdnaxeli/wishlister/pkg/repository"
)

func (a *app) CreateWishList(
	ctx context.Context,
	params CreateWishlistParams,
) (listID string, adminID string, err error) {
	if params.Name == "" {
		return "", "", ErrWishListNameEmpty
	}

	if params.Username == "" {
		return "", "", ErrWishListUsernameEmpty
	}

	if params.UserEmail == "" {
		return "", "", ErrWishListUserEmailEmpty
	}

	listID, _ = nanoid.New()
	adminID, _ = nanoid.New()

	userID, err := a.GetOrCreateUser(
		ctx,
		params.Username,
		params.UserEmail,
	)
	if err != nil {
		return "", "", err
	}

	err = a.createWishList(ctx, listID, adminID, params.Name, userID)
	if err != nil {
		return "", "", err
	}

	err = a.emailSender.SendNewWishListEmail(params.UserEmail, params.Username, listID, adminID)
	if err != nil {
		log.Print(err)
	}

	return listID, adminID, nil
}

func (a *app) createWishList(
	ctx context.Context,
	listID string,
	adminID string,
	name string,
	userID string,
) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	qtx := a.queries.WithTx(tx)
	err = qtx.CreateWishList(ctx, repository.CreateWishListParams{
		ID:      listID,
		AdminID: adminID,
		Name:    name,
		UserID:  userID,
	})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
