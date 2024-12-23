package controllers

import (
	"context"
	"ecommerce/database"
	"ecommerce/models"
	"ecommerce/utils"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

// Struct Method
// func (reciever reciever_type) AddToCart() gin.HandlerFunc {}
// Value reciever and Pointer reciever
// If we matched the Application struct data-type with method recievers data-type so we can call the function or method directly.

// Encapsulation
// Allows you to encapsulate data and methods together, similar to classes in object-oriented programming.

// Reusability
// Methods with receivers can be reused across multiple instances of a type.

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("productId")
		if productQueryId == "" {
			log.Println("Product id is empty")
			//  c.Abort -->> Simply abort to the unauthorized person like hacker without sending to the message or status code
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryId := c.Query("userId") // /route/2/:id  => 1, 2, 3
		if userQueryId == "" {
			log.Println("User id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		// We are converting the id string into ObjectID primitive form
		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productId, userQueryId)
		if err != nil {
			utils.ErrorHandler(err, c, http.StatusInternalServerError, false, err)
		}

		utils.ResponseHandler(c, 200, true, "Successfully added to the cart", nil)
		ctx.Done()
	}
}

func (app *Application) RemoveItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("productId")
		if productQueryId == "" {
			log.Println("Product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := c.Query("userId")
		if userQueryID == "" {
			log.Println("User id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productId, userQueryID)
		if err != nil {
			utils.ErrorHandler(err, c, http.StatusInternalServerError, false, err)
		}

		utils.ResponseHandler(c, 200, true, "Successfully remove item from the cart", nil)
		ctx.Done()
	}

}

func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" {
			c.JSON(http.StatusMethodNotAllowed, "Request method is invalid")
			return
		}

		user_id := c.Query("user_id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, "Invalid id")
			c.Abort()
			return
		}

		// Convert the string id into primitive object id
		userId, _ := primitive.ObjectIDFromHex(user_id)
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledCart models.User

		find := bson.D{"_id", userId}
		err := UserCollection.FindOne(ctx, find).Decode(&filledCart)
		utils.ErrorHandler(err, c, http.StatusInternalServerError, false, "Something went wrong. Please try again !")

		filter_match_stage := bson.D{"$match", bson.D{"_id", userId}}
		unwind_stage := bson.D{"$unwind", bson.D{"path", "$user_cart"}}
		// The aggregation pipeline uses the $group stage to group the documents by the _id field,
		grouping_stage := bson.D{"$group", bson.D{"_id", "$_id"}, bson.D{"total", bson.D{"$sum", "$user_cart.price"}}}

		cursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match_stage, unwind_stage, grouping_stage})
		utils.ErrorHandler(err, c, http.StatusInternalServerError, false, "Something went wrong !")

		var listing []bson.M

		for cursor.Next(ctx) {
			if err := cursor.Decode(&listing); err != nil {
				log.Println("Error while doing aggregation in GetItemFromCart function ", err)
				log.Panic(err)
			}
		}

		if err = cursor.Err(); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			log.Fatal()
		}

		defer cursor.Close(ctx)

		for _, json := range listing {
			utils.ResponseHandler(c, http.StatusOK, true, "", json["total"])
			utils.ResponseHandler(c, http.StatusOK, true, "", filledCart.User_Cart)
		}

		ctx.Done()
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Panicln("User id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryId)
		utils.ErrorHandler(err, c, http.StatusInternalServerError, false, err.Error())

		utils.ResponseHandler(c, http.StatusCreated, true, "Successfully placed the order", nil)
		ctx.Done()
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("productId")
		if productQueryId == "" {
			log.Println("Product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := c.Query("userId")
		if userQueryID == "" {
			log.Println("User id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productId, userQueryID)
		if err != nil {
			utils.ErrorHandler(err, c, http.StatusInternalServerError, false, err)
		}

		utils.ResponseHandler(c, 200, true, "Successfully placed the order", nil)
		ctx.Done()
	}
}
