package main

import (
	"ecommerce/controllers"
	"ecommerce/database"
	"ecommerce/middleware"
	"ecommerce/routes"
	"log"

	"github.com/gin-gonic/gin"

	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	// Product Data from Product Collection and User Data from User Collection
	app := controllers.NewApplication(database.ProductData(database.Client, "Product"), database.UserProduct(database.Client, "User"))

	router := gin.New() // New returns a new blank Engine instance without any middleware attached.

	router.Use(gin.Logger()) // Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port)) // when critical errors encounter in the program which stops the continuation of the program so we have to log the error messages and then immediately terminates the program with a non-zero exit status code
}
