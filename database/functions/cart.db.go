package database

import (
	"context"
	"ecommerce/models"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

	userId, err := primitive.ObjectIDFromHex(userQueryID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{{Key: "_id", Value: userId}}
	update := bson.M{"$pull": bson.M{"user_cart": bson.M{"_id": productId}}} // Removes specific elements from an array that match a condition.
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantRemoveItemFromCart
	}

	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userQueryId string) error {

	userId, err := primitive.ObjectIDFromHex(userQueryId)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var getCartItems models.User
	var orderCart models.Order

	// Making an order information for user
	orderCart.Order_ID = primitive.NewObjectID()
	orderCart.Ordered_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderCart.Order_Cart = make([]models.ProductUser, 0)
	orderCart.Payment_Method.COD = true

	// Aggregate Total Price of Items in User Cart
	// path tell us where do you want to unwind either array or slice
	unwind := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$user_cart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "$_id"},
		{Key: "total", Value: bson.D{
			{Key: "$sum", Value: "$user_cart.price"},
		}},
	},
	}}

	cursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	var allItemsAddition []bson.M

	if err = cursor.All(ctx, &allItemsAddition); err != nil {
		panic(err)
	}

	if err = cursor.Err(); err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}
	defer cursor.Close(ctx)

	var total_price int32

	for _, user_item := range allItemsAddition {
		price := user_item["total"]
		// Below uses a type assertion to convert the price value (which is user_item["total"]) to the type int32.
		// If price is not of type int32, this assertion will cause a runtime panic.
		total_price = price.(int32) // Price we got, converted to int32
	}

	orderCart.Price = int(total_price)

	// Insert the Order into the User's order list
	filter := bson.D{{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "orders", Value: orderCart}}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	// Fetch the cart items from the user collection ; Retrieves the current state of the user's cart purpose to get User_Cart items.
	find := bson.D{{Key: "_id", Value: userId}}
	err = userCollection.FindOne(ctx, find).Decode(&getCartItems)
	if err != nil {
		log.Println(err)
		return ErrCantGetItem
	}

	// Add Cart Items (which retrieves from getCartItems) to the Order's Cart
	filter2 := bson.D{{Key: "_id", Value: userId}}
	// 	Use orders.$[].order_list to update all orders.
	// Use orders.$[<identifier>].order_list with arrayFilters for selective updates.
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": getCartItems.User_Cart}}

	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}

	//  Deleting / Updating / Clear the User's Cart after transferng user cart items to the order's cart
	// Empty the UserCart after completing the purchase.
	usercart_empty := make([]models.ProductUser, 0)
	filter3 := bson.D{{Key: "_id", Value: userId}}
	update3 := bson.D{{Key: "$set", Value: bson.D{{Key: "user_cart", Value: usercart_empty}}}}

	_, err = userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	return nil
}

func InstantBuyer(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, productId primitive.ObjectID, userQueryID string) error {

	userId, err := primitive.ObjectIDFromHex(userQueryID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var product_details models.ProductUser
	var orders_detail models.Order

	orders_detail.Order_ID = primitive.NewObjectID()
	orders_detail.Ordered_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orders_detail.Order_Cart = make([]models.ProductUser, 0)
	orders_detail.Payment_Method.COD = true

	// Find that specific product by the id which user want to add them into their cart
	err = prodCollection.FindOne(ctx, bson.D{{Key: "_id", Value: productId}}).Decode(&product_details)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	// Add Orders Details into the usercollection order's
	filter := bson.D{{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "orders", Value: orders_detail}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}

	// Add Product Details into the User's Order.Order_Cart
	filter2 := bson.D{{Key: "_id", Value: userId}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": product_details}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}

	return nil
}
