-- name: CreateWishList :exec
insert into wishlists (
    id, admin_id, name, group_id
)
values (
    ?, ?, ?, ?
);
