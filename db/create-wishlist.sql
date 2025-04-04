-- name: CreateWishList :exec
insert into wishlists (
    id, admin_id, name
)
values (
    ?, ?, ?
);
