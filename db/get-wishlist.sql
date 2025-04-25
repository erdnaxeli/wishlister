-- name: GetWishList :one
select
    id,
    admin_id,
    group_id,
    name
from wishlists
where id = ?;
