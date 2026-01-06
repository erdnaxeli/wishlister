-- name: GetUserWishLists :many
select id, admin_id, name
from wishlists
where user_id = ?;
