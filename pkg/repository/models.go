// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repository

import (
	"database/sql"
)

type User struct {
	ID    string
	Name  string
	Email sql.NullString
}

type Wishlist struct {
	ID      string
	AdminID string
	UserID  sql.NullString
	Name    string
}

type WishlistElement struct {
	ID          string
	WishlistID  string
	Name        string
	Description sql.NullString
	Url         sql.NullString
}
