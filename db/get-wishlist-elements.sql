-- name: GetWishListElements :many
select
    id,
    name,
    description,
    url
from wishlist_elements
where wishlist_id = ?;
