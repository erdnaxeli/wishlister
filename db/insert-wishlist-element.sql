-- name: InsertWishListElement :exec
insert into wishlist_elements (
    id,
    wishlist_id,
    name,
    description,
    url
) values (
    ?, ?, ?, ?, ?
);
