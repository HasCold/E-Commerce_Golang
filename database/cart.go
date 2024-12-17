package controllers

import "errors"

// Custom Error Defining :- when we are interacting with the DB perform some operation and cause error
var (
	ErrCantFindProduct        = errors.New("can't find the product")
	ErrCantDecodeProducts     = errors.New("can't find the products")
	ErrUserIdIsNotValid       = errors.New("this user is not valid")
	ErrCantUpdateUser         = errors.New("cannot add this product to the cart")
	ErrCantRemoveItemFromCart = errors.New("cannot remove this item from the cart")
	ErrCantGetItem            = errors.New("was unable to get the items from the cart")
	ErrCantBuyCartItem        = errors.New("cannot update the purchase")
)

func AddProductToCart() {

}

func RemoveCartItem() {

}

func BuyItemFromCart() {

}

func InstantBuyer() {

}
