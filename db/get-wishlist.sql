-- name: GetWishList :one
select
    id,
    admin_id,
    name
from wishlists
where id = ?;
