// Package wishlister implements a wishlists manager.
package wishlister

import (
	"context"
	"database/sql"
	"fmt"

	nanoid "github.com/matoous/go-nanoid/v2"

	"github.com/erdnaxeli/wishlister/pkg/email"
	"github.com/erdnaxeli/wishlister/pkg/repository"

	// register sqlite driver
	_ "modernc.org/sqlite"
)

// CreateWishlistParams represents the parameters to create a new wishlist.
type CreateWishlistParams struct {
	Name      string
	Username  string
	UserEmail string
}

// CreateGroupParams represents the parameters to create a new group.
type CreateGroupParams struct {
	Name      string
	UserEmail string
}

// Session represents a user session.
type Session struct {
	UserID         string
	Username       string
	UserEmail      string
	SessionID      string
	MagicLinkToken string
}

// WishList represents a wishlist.
type WishList struct {
	ID string

	Name     string
	AdminID  string
	GroupID  string
	Username string

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
	CreateGroup(ctx context.Context, params CreateGroupParams) (string, error)

	// Create a new wishlist.
	//
	// Return the wishlist id and the admin id.
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
	// If the wishlist is not found, an error ErrWishListNotFound is returned.
	// If the adminId token is incorrect, an error ErrWishListInvalidAdminId is returned.
	GetEditableWishList(ctx context.Context, listID string, adminID string) (WishList, error)

	// GetUserWishLists returns all wishlists for a given user.
	//
	// Elements are not included in the returned wishlists.
	GetUserWishLists(ctx context.Context, userID string) ([]WishList, error)

	// UpdateListElements updates the elements of a wishlist.
	//
	// This method check that the adminId token is the correct one for this wishlist.
	//
	// If the wishlist is not found, an error ErrWishListNotFound is returned.
	// If the adminId token is incorrect, an error ErrWishListInvalidAdminId is returned.
	//
	// The elements parameter is the full list of elements to set on the wishlist.
	// Existing elements are deleted and replaced by the new ones.
	UpdateListElements(
		ctx context.Context,
		listID string,
		adminID string,
		elements []WishListElement,
	) error

	// SendMagicLink sends a magic link to the given email address.
	//
	// The link can be used to login the user.
	SendMagicLink(ctx context.Context, email string) error

	// GetSession returns the user id associated with the given session id.
	//
	// If the session is not found, an error ErrSessionNotFound is returned.
	GetSession(ctx context.Context, sessionID string) (Session, error)

	// GetSessionByMagicLink returns the user id associated with the given magic link token.
	//
	// If the session is not found, an error ErrSessionNotFound is returned.
	// Once used, the magic link token is invalidated and cannot be used again.
	GetSessionByMagicLink(ctx context.Context, token string) (Session, error)
}

type app struct {
	db      *sql.DB
	queries *repository.Queries

	emailSender email.Sender
}

// New returns an App object with the default config.
func New(emailSender email.Sender) (App, error) {
	return NewWithConfig("db.sqlite", emailSender)
}

// NewWithConfig returns an App object with a specific config.
//
// dbFile is the path to the sqlite db file.
func NewWithConfig(dbFile string, emailSender email.Sender) (App, error) {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("error while opening database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error while pinging database: %w", err)
	}

	return &app{
		db:      db,
		queries: repository.New(db),

		emailSender: emailSender,
	}, nil
}

func (a *app) CreateGroup(ctx context.Context, params CreateGroupParams) (string, error) {
	groupID, _ := nanoid.New()

	err := a.queries.CreateGroup(
		ctx,
		repository.CreateGroupParams{
			ID:   groupID,
			Name: params.Name,
		},
	)
	if err != nil {
		return "", err
	}

	return groupID, nil
}

func (a *app) GetGroup(_ context.Context, _ string) {}

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
