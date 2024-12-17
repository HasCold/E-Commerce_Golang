package controllers

import (
	"context"
	"net/http"
	"time"

	"ecommerce/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func HashPassword(password string) string {

}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {

}

func Signup() gin.HandlerFunc {

	// Closure Func
	return func(c *gin.Context) { // c is a pointer to a gin.Context struct
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second) // Timeout Handling: Ensures that long-running operations (like database queries or API calls) don't block your application indefinitely.
		defer cancel()                                                            //  defer will execute this at the end of all nearby function execution ; Resource Management.

		if c.Request.Method != "POST" {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Request method is invalid"})
		}

		var user model.User

		// Binds JSON data from the request body to the user variable. If the request body contains invalid JSON or doesn't match the User struct, it returns an error.
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if validationErr := validate.Struct(user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
		}
	}

}

func Login() gin.HandlerFunc {

}
