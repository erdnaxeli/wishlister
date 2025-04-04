package wishlister

// WishListNotFoundError is the error when a wishlist cannot be found.
type WishListNotFoundError struct{}

func (WishListNotFoundError) Error() string {
	return "no wishlist found"
}

// WishListInvalidAdminIDError is the error when a given adminID token does not
// match the given wishlist (or listID).
type WishListInvalidAdminIDError struct{}

func (WishListInvalidAdminIDError) Error() string {
	return "access denied to this wishlist"
}
