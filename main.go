package main

import (
	"ecommerce/constants"
	"ecommerce/controllers"
	"ecommerce/database"
	"ecommerce/middleware"
	"ecommerce/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var port string

// Init() will execute before the main funcion or main goroutine.
func init() {
	if _, err := os.Stat(".env"); err == nil {
		// Load environment variables from .env file
		constants.LoadENV()
		port = constants.PORT
		if port == "" {
			port = "8000"
		}
	}

}

func main() {

	// Product Data from Product Collection and User Data from User Collection
	// Cart Controller
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.Default() // Default returns a gin engine instance which is used to build a middleware, logger and routing purposes. creates a new Gin router with two middlewares already included : Logger and Recovery Middleware

	routes.UserRoutes(router)

	// Pass the middleware in Use method
	router.Use(middleware.Authentication())
	// Below are the api's will authorize first from the middleware
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItemFromCart())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	router.GET("/listcart", controllers.GetItemFromCart())
	router.POST("/addaddress", controllers.AddAddress())
	router.PUT("/edithomeaddress", controllers.EditHomeAddress())
	router.PUT("/editworkaddress", controllers.EditWorkAddress())
	router.GET("/deleteaddresses", controllers.DeleteAddress())

	log.Fatal(router.Run(":" + port)) // when critical errors encounter in the program which stops the continuation of the program so we have to log the error messages and then immediately terminates the program with a non-zero exit status code
}
