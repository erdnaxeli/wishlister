package wishlister

type WishListNotFoundError struct{}

func (WishListNotFoundError) Error() string {
	return "no wishlist found"
}

type WishListInvalidAdminIdError struct{}

func (WishListInvalidAdminIdError) Error() string {
	return "access denied to this wishlist"
}
