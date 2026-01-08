package wishlister

import "errors"

// ErrWishListNotFound is the error when a wishlist cannot be found.
var ErrWishListNotFound = errors.New("no wishlist found")

// ErrWishListInvalidAdminID is the error when a given adminID token does not
// match the given wishlist (or listID).
var ErrWishListInvalidAdminID = errors.New("access denied to this wishlist")

// ErrWishListNameEmpty is returned when the wishlist name is empty.
var ErrWishListNameEmpty = errors.New("wishlist name cannot be empty")

// ErrWishListUsernameEmpty is returned when the wishlist username is empty.
var ErrWishListUsernameEmpty = errors.New("wishlist username cannot be empty")

// ErrSessionNotFound is returned when a session cannot be found.
var ErrSessionNotFound = errors.New("session not found")
