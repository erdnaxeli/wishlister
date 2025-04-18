// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: create-wishlist.sql

package repository

import (
	"context"
)

const createWishList = `-- name: CreateWishList :exec
insert into wishlists (
    id, admin_id, name
)
values (
    ?, ?, ?
)
`

type CreateWishListParams struct {
	ID      string
	AdminID string
	Name    string
}

func (q *Queries) CreateWishList(ctx context.Context, arg CreateWishListParams) error {
	_, err := q.db.ExecContext(ctx, createWishList, arg.ID, arg.AdminID, arg.Name)
	return err
}
