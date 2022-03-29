package database

import "errors"

var (
	ErrCantFindProduct    = errors.New("can't find the product")
	ErrCantDecodeProduct  = errors.New("can't decode the product")
	ErrUserIdIsNotValid   = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("cannot add this product to the cart")
	ErrCantRemoveItemCart = errors.New("cannot remove this item from the cart")
	ErrCantGetItem        = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
)

/**
* #Add Product To Cart Function
*
*
**/
func AddProductToCart() {

}

/**
* #Remove Cart Item Function
*
*
**/
func RemoveCartItem() {

}

/**
* #Buy Item From Cart Function
*
*
**/
func BuyItemFromCart() {

}

/**
* #Instant Buy Function
*
*
**/
func InstantBuy() {

}