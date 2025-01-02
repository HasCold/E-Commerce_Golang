package controllers

import (
	"context"
	"ecommerce/models"
	"ecommerce/utils"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func User_Test() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request method is invalid !"})
			return
		}

		user_id := c.Param("userId")

		if user_id == "" {
			log.Println("User Id is empty")
			c.AbortWithError(http.StatusNotFound, errors.New("something missing "))
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id !"})
			return
		}

		userData := bson.M{}
		filter := bson.D{{Key: "_id", Value: userId}}
		// filter := bson.M{}
		err = UserCollection.FindOne(ctx, filter).Decode(&userData)
		if err != nil {
			log.Println(err)
			utils.ErrorHandler(c, http.StatusInternalServerError, false, err.Error())
			return
		}

		utils.ResponseHandler(c, http.StatusOK, true, "", userData)
		ctx.Done()
	}
}

func All_User_Test() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request method is invalid !"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var userData []models.User

		opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})
		cursor, err := UserCollection.Find(ctx, bson.D{}, opts)
		if err != nil {
			log.Println(err)
			utils.ErrorHandler(c, http.StatusInternalServerError, false, err.Error())
			return
		}

		for cursor.Next(ctx) {
			var user models.User
			if err := cursor.Decode(&user); err != nil {
				log.Println(err)
				log.Panic()
			}
			userData = append(userData, user)
		}

		if err = cursor.Err(); err != nil {
			log.Println("Error while decoding the user")
			c.Abort()
			return
		}
		defer cursor.Close(ctx)

		utils.ResponseHandler(c, http.StatusOK, true, "All Users Data", userData)
		ctx.Done()
	}
}

func Test_Empty_Order_Cart() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request method is invalid !"})
			return
		}

		user_id := c.Param("userId")

		if user_id == "" {
			log.Println("User Id is empty")
			c.AbortWithError(http.StatusNotFound, errors.New("something missing "))
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id !"})
			return
		}

		Order := make([]models.Order, 0)

		filter := bson.D{{Key: "_id", Value: userId}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "orders", Value: Order}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err)
			utils.ErrorHandler(c, http.StatusInternalServerError, false, err.Error())
			return
		}

		utils.ResponseHandler(c, http.StatusOK, true, "Order's Cart Successfully Empty", nil)
		ctx.Done()
	}
}
