package controllers

import (
	"context"
	"net/http"
	"time"

	"ecommerce/model"
	"ecommerce/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		utils.ErrorHandler(err, c, http.StatusInternalServerError, false, err)

		if count > 0 {
			utils.ResponseHandler(c, http.StatusBadRequest, false, "User already exits")
		}

		count, err = UserCollection.CountDocument(ctx, bson.M{"phone": user.phone})
		utils.ErrorHandler(err, c, http.StatusBadRequest, false, err)

		if count > 0 {
			utils.ResponseHandler(c, http.StatusBadRequest, false, "Phone number already exists")
		}

		password := HashPassword(*user.Password)
		user.Password = &password // rather than copying the whole password we have to pass the actual memory address reference of password

		// RFC3339 is a standard date-time format, such as:  2024-12-18T14:30:15Z. To ensure the things are consistent in DB
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID
		user.User_ID = user.ID.Hex() // Hex returns the hex encoding of the ObjectID as a string.

		token, refresh_token, _ := utils.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, *user.User_ID)

		user.Token = &token
		user.Refresh_Token = &refresh_token
		user.User_Cart = make([]model.ProductUser, 0) // make a built-in datatype or function helps to create and initialize different data-types like slices, map and channels
		user.Address_Details = make([]model.Address, 0)
		user.Order_Status = make([]model.Order, 0)

		_, insertErr := UserCollection.InsertOne(ctx, user)
		utils.ErrorHandler(insertErr, c, http.StatusInternalServerError, false, insertErr)

		utils.ResponseHandler(c, http.StatusCreated, true, "Successfully Signed In !")
	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()

		if c.Request.Method != "POST" {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Request Method is Invalid !"})
		}

		var user model.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

	}
}

// bson.M{}  -->>  M is an unordered representation of a BSON document. This type should be used when the order of the elements does not matter. This type is handled as a regular map[string]interface{} when encoding and decoding.
// bson.D{}  -->> D is an ordered representation of a BSON document. This type should be used when the order of the elements matters, such as MongoDB command documents.
