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
	// Cart Controller
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.Default() // Default returns a gin engine instance which is used to build a middleware, logger and routing purposes. creates a new Gin router with two middlewares already included : Logger and Recovery Middleware

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItemFromCart())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port)) // when critical errors encounter in the program which stops the continuation of the program so we have to log the error messages and then immediately terminates the program with a non-zero exit status code
}
