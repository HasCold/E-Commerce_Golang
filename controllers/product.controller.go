package controllers

import (
	"context"
	"ecommerce/database"
	"ecommerce/model"
	"ecommerce/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ProdCollection *mongo.Collection = database.ProductData(database.Client, "Products")

func GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method != "POST" {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Request method is invalid !"})
			return
		}

		var ProductList []model.Product

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Each section uses the following cursor variable, which is a Cursor struct that contains all the documents in a collection:
		findOption := options.Find()
		cursor, err := ProdCollection.Find(ctx, bson.D{}, findOption)
		utils.ErrorHandler(err, c, http.StatusInternalServerError, false, "Something went wrong. Please try again !")

		// To retrieve documents from your cursor individually while blocking the current goroutine, use the Next() method.
		for cursor.Next(ctx) {
			err := cursor.Decode(&ProductList) // Bind the individual product within the productList
			if err != nil {
				log.Fatal(err)
			}
		}

		if err := cursor.Err(); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			log.Fatal(err)
		}

		defer cursor.Close(ctx)

		utils.ResponseHandler(c, 200, true, "Successfully get all the products", ProductList)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Request method is invalid !"})
			return
		}

		var searchProducts []model.Product

		queryParam := c.Query("name") // ?name="iphone"
		if queryParam == "" {
			log.Println("query name is empty")
			c.Header("Content-Type", "application/json")
			c.Abort() // -->>This can be used when you need to stop further processing due to invalid input, unauthorized access, etc.
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		findOptions := options.Find()
		search := bson.M{
			"product_name": bson.M{
				"$regex": queryParam,
			},
		}
		cursor, err := ProdCollection.Find(ctx, search, findOptions)
		utils.ErrorHandler(err, c, http.StatusNotFound, false, "Something went wrong. Please try again")

		// To retrieve documents from your cursor individually while blocking the current goroutine, use the Next() method.
		for cursor.Next(ctx) {
			if err := cursor.Decode(&searchProducts); err != nil {
				log.Fatal(err)
			}
		}

		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, "Invalid request")
			log.Fatal(err)
		}

		defer cursor.Close(ctx)

		utils.ResponseHandler(c, http.StatusOK, true, "", searchProducts)
	}
}

func ProductViewerAdmin() gin.HandlerFunc {

}

// Match criteria with literal values use the following format:
// filter := bson.D{{"<field>", "<value>"}}

// Match criteria with a query operator use the following format:
// filter := bson.D{{"<field>", bson.D{{"<operator>", "<value>"}}}}
