package wishlister

import (
	"context"
	"database/sql"
	"errors"
)

func (a *app) GetWishList(ctx context.Context, listID string) (WishList, error) {
	wishList, err := a.getWishList(ctx, listID)
	if err != nil {
		return WishList{}, err
	}

	// hide admin ID
	wishList.AdminID = ""
	err = a.populateElements(ctx, &wishList)
	if err != nil {
		return WishList{}, err
	}

	return wishList, nil
}

func (a *app) GetEditableWishList(
	ctx context.Context,
	listID string,
	adminID string,
) (WishList, error) {
	list, err := a.checkListEditAccess(ctx, listID, adminID)
	if err != nil {
		return WishList{}, err
	}

	err = a.populateElements(ctx, &list)
	if err != nil {
		return WishList{}, err
	}

	return list, nil
}

func (a *app) checkListEditAccess(
	ctx context.Context,
	listID string,
	adminID string,
) (WishList, error) {
	wishList, err := a.getWishList(ctx, listID)
	if err != nil {
		return WishList{}, err
	}

	if wishList.AdminID != adminID {
		return WishList{}, ErrWishListInvalidAdminID
	}

	return wishList, nil
}

func (a *app) getWishList(ctx context.Context, listID string) (WishList, error) {
	list, err := a.queries.GetWishList(ctx, listID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return WishList{}, ErrWishListNotFound
		}
		return WishList{}, err
	}

	wishList := WishList{
		AdminID:  list.AdminID,
		ID:       list.ID,
		Name:     list.Name,
		GroupID:  list.GroupID.String,
		Username: list.Username,
	}
	return wishList, nil
}
