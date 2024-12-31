package controllers

import (
	"context"
	"ecommerce/database"
	"ecommerce/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ProdCollection *mongo.Collection = database.ProductData(database.Client, "Products")

func GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method != "GET" {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Request method is invalid !"})
			return
		}

		var ProductList []models.Product

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Each section uses the following cursor variable, which is a Cursor struct that contains all the documents in a collection:
		cursor, err := ProdCollection.Find(ctx, bson.D{})
		if err != nil {
			log.Println(err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong. Please try again !",
			})
			return
		}

		// To retrieve documents from your cursor individually while blocking the current goroutine, use the Next() method.
		for cursor.Next(ctx) {
			var product models.Product
			err := cursor.Decode(&product) // Bind the individual product within the productList
			if err != nil {
				log.Fatal(err)
			}

			ProductList = append(ProductList, product)
		}

		if err := cursor.Err(); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			log.Fatal(err)
		}

		defer cursor.Close(ctx)

		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"message":  "Successfully get all the products",
			"products": ProductList,
		})
		ctx.Done()
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Request method is invalid !"})
			return
		}

		var searchProducts []models.Product

		queryParam := c.Query("name") // ?name="iphone"
		if queryParam == "" {
			log.Println("query name is empty")
			c.Header("Content-Type", "application/json")
			c.Abort() // -->>This can be used when you need to stop further processing due to invalid input, unauthorized access, etc.
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// The $regex operator is used to perform a regex search with the case-insensitive(Mean uppercase and lowercase letter treated be as same) option $options: 'i', and the caret ^ in the regex pattern ensures that the search starts with the given letters.
		search := bson.M{
			"product_name": bson.M{
				"$regex":   "^" + queryParam,
				"$options": "i",
			},
		}
		cursor, err := ProdCollection.Find(ctx, search)
		if err != nil {
			log.Println(err)

			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Something went wrong. Please try again !",
			})
			return
		}

		// To retrieve documents from your cursor individually while blocking the current goroutine, use the Next() method.
		for cursor.Next(ctx) {
			var search models.Product
			if err := cursor.Decode(&search); err != nil {
				log.Fatal(err)
			}
			searchProducts = append(searchProducts, search)
		}

		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, "Invalid request")
			log.Fatal(err)
		}

		defer cursor.Close(ctx)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    searchProducts,
		})
		// utils.ResponseHandler(c, http.StatusOK, true, "", searchProducts)
		ctx.Done()
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method != "POST" {
			c.JSON(http.StatusBadRequest, "Request method is invalid")
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var products models.Product
		err := c.BindJSON(&products)
		if err != nil {
			log.Println(err)

			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		products.Product_ID = primitive.NewObjectID()
		_, err = ProdCollection.InsertOne(ctx, products)
		if err != nil {
			log.Println(err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Successfully added our Product Admin!!",
		})
		ctx.Done()
	}
}

// Match criteria with literal values use the following format:
// filter := bson.D{{"<field>", "<value>"}}

// Match criteria with a query operator use the following format:
// filter := bson.D{{"<field>", bson.D{{"<operator>", "<value>"}}}}
