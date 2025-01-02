package routes

import (
	"ecommerce/controllers"

	"github.com/gin-gonic/gin"
)

func TestRoutes(incomingRequest *gin.Engine) {
	incomingRequest.GET("/testuser/:userId", controllers.User_Test())
	incomingRequest.GET("/allusertest", controllers.All_User_Test())
	incomingRequest.GET("/emptyordercart/:userId", controllers.Test_Empty_Order_Cart())
}
