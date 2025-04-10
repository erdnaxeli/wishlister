// Package wishlister implements a wishlists manager.
package wishlister

import (
	"context"
	"database/sql"
	"errors"

	"github.com/erdnaxeli/wishlister/pkg/repository"
	nanoid "github.com/matoous/go-nanoid/v2"

	// register sqlite driver
	_ "modernc.org/sqlite"
)

// CreateWishlistParams represents the parameters to create a new wishlist.
type CreateWishlistParams struct {
	Name string
}

// WishList represents a wishlist.
type WishList struct {
	ID string

	Name    string
	AdminID string

	Elements []WishListElement
}

// WishListElement represents a wishlist element.
type WishListElement struct {
	Name        string
	Description string
	URL         string
}

// App is the main interface of this package.
//
// It implements all method to manage wishlists.
type App interface {
	// Create a new group.
	//
	// Return the group id.
	CreateGroup(ctx context.Context) string

	// Create a new wish list.
	//
	// Return the wish list id.
	CreateWishList(ctx context.Context, params CreateWishlistParams) (string, string, error)
	GetGroup(ctx context.Context, groupID string)

	// Get a wishlist.
	//
	// The AdminId field on the returned Wishlist object will be empty.
	//
	// If the wishlist is not found, an error WishListNotFoundError is returned.
	GetWishList(ctx context.Context, listID string) (WishList, error)

	// Get a wishlist to be edited.
	//
	// This method check that the adminId token is the correct one for this wishlist.
	//
	// If the wishlist is not found, an error WishListNotFoundError is returned.
	// If the adminId token is incorrect, an error WishListInvalidAdminIdError is returned.
	GetEditableWishList(ctx context.Context, listID string, adminID string) (WishList, error)

	UpdateListElements(
		ctx context.Context,
		listID string,
		adminID string,
		elements []WishListElement,
	) error
}

type app struct {
	db      *sql.DB
	queries *repository.Queries
}

// New returns an App object with the default config.
func New() (App, error) {
	return NewWithConfig("db.sqlite")
}

// NewWithConfig returns an App object with a specific config.
//
// dbFile is the path to the sqlite db file.
func NewWithConfig(dbFile string) (App, error) {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	return &app{
		db:      db,
		queries: repository.New(db),
	}, nil
}

func (a *app) CreateGroup(_ context.Context) string {
	groupID, _ := nanoid.New()
	return groupID
}

func (a *app) CreateWishList(
	ctx context.Context,
	params CreateWishlistParams,
) (string, string, error) {
	listID, _ := nanoid.New()
	adminID, _ := nanoid.New()

	err := a.queries.CreateWishList(ctx, repository.CreateWishListParams{
		ID:      listID,
		AdminID: adminID,
		Name:    params.Name,
	})
	if err != nil {
		return "", "", err
	}

	return listID, adminID, nil
}

func (a *app) GetGroup(_ context.Context, _ string) {}

func (a *app) GetWishList(ctx context.Context, listID string) (WishList, error) {
	list, err := a.queries.GetWishList(ctx, listID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return WishList{}, WishListNotFoundError{}
		}
		return WishList{}, err
	}

	wishList := WishList{
		ID:   list.ID,
		Name: list.Name,
	}

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

func (a *app) UpdateListElements(
	ctx context.Context,
	listID string,
	adminID string,
	elements []WishListElement,
) (err error) {
	_, err = a.checkListEditAccess(ctx, listID, adminID)
	if err != nil {
		return err
	}

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
	err = qtx.DeleteWishListElements(ctx, listID)
	if err != nil {
		return err
	}

	for _, element := range elements {
		elementID, _ := nanoid.New()
		err = qtx.InsertWishListElement(
			ctx,
			repository.InsertWishListElementParams{
				ID:          elementID,
				WishlistID:  listID,
				Name:        element.Name,
				Description: NewNullString(element.Description),
				Url:         NewNullString(element.URL),
			},
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (a *app) checkListEditAccess(
	ctx context.Context,
	listID string,
	adminID string,
) (WishList, error) {
	list, err := a.queries.GetWishList(ctx, listID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return WishList{}, WishListNotFoundError{}
		}
		return WishList{}, err
	}

	if list.AdminID != adminID {
		return WishList{}, WishListInvalidAdminIDError{}
	}

	return WishList{
		ID:      list.ID,
		AdminID: list.AdminID,
		Name:    list.Name,
	}, nil
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

// NewNullString convert a string value to a sql.NullString value.
//
// If the string is empty, the NullString is invalid, else it is valid and contains
// the string value.
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
