package database

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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

// Database Level Function

func AddProductToCart(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, productId primitive.ObjectID, userQueryID string) error {

}

func RemoveCartItem(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, productId primitive.ObjectID, userQueryID string) error {

}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userQueryId string) error {

}

func InstantBuyer(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, productId primitive.ObjectID, userQueryID string) error {

}
