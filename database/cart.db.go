package database

import (
	"context"
	"ecommerce/models"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
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

	find := bson.M{"_id": productId}
	searchFromDB, err := prodCollection.Find(ctx, find)
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var productCart []models.ProductUser

	for searchFromDB.Next(ctx) {
		if err := searchFromDB.Decode(&productCart); err != nil {
			log.Println(err)
			return ErrCantDecodeProducts
		}
	}

	if err = searchFromDB.Err(); err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}
	defer searchFromDB.Close(ctx)

	userId, err := primitive.ObjectIDFromHex(userQueryID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "user_cart", Value: productCart}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}

	return nil
}

func RemoveCartItem(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, productId primitive.ObjectID, userQueryID string) error {

}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userQueryId string) error {

}

func InstantBuyer(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, productId primitive.ObjectID, userQueryID string) error {

}
