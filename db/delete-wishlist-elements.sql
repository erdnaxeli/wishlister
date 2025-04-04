-- name: DeleteWishListElements :exec
delete from wishlist_elements
where wishlist_id = ?;
