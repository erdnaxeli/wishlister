// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: get-wishlist.sql

package repository

import (
	"context"
)

const getWishList = `-- name: GetWishList :one
select
    id,
    admin_id,
    name
from wishlists
where id = ?
`

type GetWishListRow struct {
	ID      string
	AdminID string
	Name    string
}

func (q *Queries) GetWishList(ctx context.Context, id string) (GetWishListRow, error) {
	row := q.db.QueryRowContext(ctx, getWishList, id)
	var i GetWishListRow
	err := row.Scan(&i.ID, &i.AdminID, &i.Name)
	return i, err
}
