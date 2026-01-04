-- name: GetWishList :one
select
    wishlists.id,
    admin_id,
    group_id,
    wishlists.name,
    users.name as username
from wishlists
join users on wishlists.user_id = users.id
where wishlists.id = ?;
