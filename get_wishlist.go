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

func (a *app) GetUserWishLists(ctx context.Context, userID string) ([]WishList, error) {
	listsData, err := a.queries.GetUserWishLists(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []WishList{}, nil
		}

		return nil, err
	}

	var wishLists []WishList
	for _, listData := range listsData {
		wishLists = append(wishLists, WishList{
			ID:      listData.ID,
			AdminID: listData.AdminID,
			Name:    listData.Name,
		})
	}

	return wishLists, nil
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

func (a *app) populateElements(ctx context.Context, list *WishList) error {
	elements, err := a.queries.GetWishListElements(ctx, list.ID)
	if err != nil {
		return err
	}

	for _, element := range elements {
		list.Elements = append(
			list.Elements,
			WishListElement{
				Name:        element.Name,
				Description: element.Description.String,
				URL:         element.Url.String,
			},
		)
	}

	return nil
}
