package routes

import (
	"ecommerce/controllers"

	"github.com/gin-gonic/gin"
)

// A pointer is a variable that can store the actual memory address locaation of another variable.
// gin.Engine is the main struct that represents the HTTP router and serves as the foundation of a Gin application. It manages the routing of incoming HTTP requests to the appropriate handlers and provides middleware support.

func UserRoutes(incomingRequest *gin.Engine) { // incomingRequest is a pointer to a gin.Engine struct
	incomingRequest.POST("/users/signup", controllers.SignUp())
	incomingRequest.POST("/users/login", controllers.Login())
	incomingRequest.POST("/admin/addproduct", controllers.ProductViewerAdmin())

	incomingRequest.POST("/users/productview", controllers.GetAllProducts())
	incomingRequest.POST("/user/search", controllers.SearchProductByQuery())
}
