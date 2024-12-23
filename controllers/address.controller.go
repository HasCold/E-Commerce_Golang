package controllers

import (
	"context"
	"ecommerce/models"
	"ecommerce/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request method"})
			return
		}

		user_id := c.Query("user_id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.Abort()
			utils.ResponseHandler(c, http.StatusNotFound, false, "Invalid Code", nil)
		}

		// Convert the string id into primitive object id
		userId, err := primitive.ObjectIDFromHex(user_id)
		utils.ErrorHandler(err, c, http.StatusInternalServerError, false, "Internal Server Error")

		var addresses models.Address

		addresses.Address_ID = primitive.NewObjectID()
		if err = c.BindJSON(&addresses); err != nil {
			c.JSON(http.StatusNotAcceptable, err.Error())
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Role Of Path :=
		// . It tells MongoDB which array field to unwind (i.e., break apart).
		// . The path is specified as a field name and must be prefixed with $ to denote the field within the document.

		match_filter_stage := bson.D{"$match", bson.D{"_id", userId}}
		unwind_stage := bson.D{"$unwind", bson.D{"path", "$address"}}
		// The documents will be grouped based on the address_id field.
		// For each group of documents that share the same address_id, the count field will be calculated as the sum of 1s for each document in the group.
		grouping_stage := bson.D{"$group", bson.D{"_id", "$address_id"}, bson.D{"count", bson.D{"$sum", 1}}}

		cursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline(match_filter_stage, unwind_stage, grouping_stage))
		utils.ErrorHandler(err, c, http.StatusInternalServerError, false, "Internal Server Error")

		var addressInfo []bson.M

		for cursor.Next(addressInfo) {
			if err := cursor.Decode(&addressInfo); err != nil {
				log.Println("Error while doing aggregation in AddAddress function ", err)
				log.Panic(err)
			}
		}

		if err = cursor.Err(); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			panic(err)
		}

		defer cursor.Close(ctx)

		var size int32
		for _, addressNo := range addressInfo {
			count := addressNo["count"] // The type of count is interface{} because Go maps with interface{} values can store values of any type.
			// If count is not of type int32, the type assertion count.(int32) will cause a runtime panic.
			size = count.(int32) //The syntax count.(int32) asserts that the underlying type of count is int32.
		}

		// If User Addresses is less than 2
		if size < 2 {
			filter := bson.D{"_id", userId}
			update := bson.D{"$push", bson.D{"address", addresses}}

			_, err := UserCollection.UpdateOne(ctx, filter, update)
			utils.ErrorHandler(err, c, http.StatusInternalServerError, false, "Something went wrong")

		} else {
			utils.ResponseHandler(c, http.StatusBadRequest, false, "Not Allowed", nil)
		}

		ctx.Done()
	}
}

func EditHomeAddress() gin.HandlerFunc {

}

func EditWorkAddress() gin.HandlerFunc {

}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Request method is Invalid"})
		}

		user_id := c.Query("user_id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Search Index is Invalid"})
			c.Abort()
			return
		}

		addresses := make([]models.Address, 0)
		// Convert the string ID into primitive Object ID
		userId, err := primitive.ObjectIDFromHex(user_id)
		utils.ErrorHandler(err, c, http.StatusInternalServerError, false, "Internal Server Error")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{"_id", bson.D{"$eq", userId}}
		update := bson.D{"$set", bson.D{"address", addresses}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		utils.ErrorHandler(err, c, http.StatusInternalServerError, false, "Something went wrong. Please Try Again !")

		utils.ResponseHandler(c, http.StatusOK, true, "Successfully Deleted", nil)
		ctx.Done()
	}
}
