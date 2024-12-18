package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"ecommerce/database"
	"ecommerce/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate *validator.Validate
var UserCollection *mongo.Collection = database.UserProduct(database.Client, "Users")

func HashPassword(password string) string {

}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {

}

func Signup() gin.HandlerFunc {

	// Closure Func
	return func(c *gin.Context) { // c is a pointer to a gin.Context struct
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second) // Timeout Handling: Ensures that long-running operations (like database queries or API calls) don't block your application indefinitely.
		defer cancel()                                                           //  defer will execute this at the end of all nearby function execution ; Resource Management.

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

		count, err := UserCollection.CountDocument(ctx, bson.M{"email": user.email})
		if err != nil {
			// Triggers a panic: After logging the error, it calls the panic() function, which stops the normal execution of the program and begins the unwinding of the stack. This allows deferred functions to execute before the program terminates.
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "User already exist",
			})
		}

		count, err = UserCollection.CountDocument(ctx, bson.M{"phone": user.phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err,
			})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Phone Number already exist",
			})
		}

	}

}

func Login() gin.HandlerFunc {

}

// bson.M{}  -->>  M is an unordered representation of a BSON document. This type should be used when the order of the elements does not matter. This type is handled as a regular map[string]interface{} when encoding and decoding.
// bson.D{}  -->> D is an ordered representation of a BSON document. This type should be used when the order of the elements matters, such as MongoDB command documents.
