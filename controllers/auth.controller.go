package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"ecommerce/config"
	"ecommerce/database"
	"ecommerce/models"
	"ecommerce/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var validate *validator.Validate = validator.New()
var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")

func HashPassword(password string) string {
	// ASCII Characters (Basic English):
	// Each character takes 1 byte (e.g., A, z, 0, !).
	// 72 ASCII characters = 72 bytes.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, hashPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(userPassword))
	valid := true
	msg := ""

	if err != nil {
		msg = "Login or Password is incorrect"
		valid = false
	}

	return valid, msg
}

func SignUp() gin.HandlerFunc {
	// Closure Func
	return func(c *gin.Context) { // c is a pointer to a gin.Context struct
		if c.Request.Method != "POST" {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Request method is invalid"})
			return
		}

		config.TokenSetting()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second) // Timeout Handling: Ensures that long-running operations (like database queries or API calls) don't block your application indefinitely.
		defer cancel()                                                            //  defer will execute this at the end of all nearby function execution ; Resource Management.

		var user models.User

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
			return
		}

		filter := bson.D{{Key: "email", Value: user.Email}}
		count, err := UserCollection.CountDocuments(ctx, filter)
		if err != nil {
			log.Println(err)
			// Triggers a panic: After logging the error, it calls the panic() function, which stops the normal execution of the program and begins the unwinding of the stack. This allows deferred functions to execute before the program terminates.
			// log.Panic(err)
			utils.ErrorHandler(c, http.StatusInternalServerError, false, err.Error())
			return
		}

		if count > 0 {
			utils.ResponseHandler(c, http.StatusBadRequest, false, "User already exits", nil)
			return
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Println(err)
			utils.ResponseHandler(c, http.StatusBadRequest, false, err.Error(), nil)
			return
		}

		if count > 0 {
			utils.ResponseHandler(c, http.StatusBadRequest, false, "Phone number already exists", nil)
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		// RFC3339 is a standard date-time format, such as:  2024-12-18T14:30:15Z. To ensure the things are consistent in DB
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		hexValue := user.ID.Hex() // Hex returns the hex encoding of the ObjectID as a string.
		user.User_ID = &hexValue

		token, refresh_token, _ := config.JwtWrapper.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, *user.User_ID)

		user.Token = &token
		user.Refresh_Token = &refresh_token
		user.User_Cart = make([]models.ProductUser, 0) // make a built-in datatype or function helps to create and initialize different data-types like slices, map and channels
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		_, insertErr := UserCollection.InsertOne(ctx, user)
		if insertErr != nil {
			log.Println(insertErr)
			utils.ErrorHandler(c, http.StatusInternalServerError, false, insertErr.Error())
			return
		}

		utils.ResponseHandler(c, http.StatusCreated, true, "Successfully Signed Up !", nil)
		ctx.Done()
	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Request Method is Invalid !"})
			return
		}

		config.TokenSetting()

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"email": user.Email} // Map
		err := UserCollection.FindOne(ctx, filter).Decode(&foundUser)
		if err != nil {
			log.Println(err)
			utils.ErrorHandler(c, http.StatusInternalServerError, false, "Login Email or Password is Incorrect !")
			return
		}

		PasswordisValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

		if !PasswordisValid {
			log.Println(msg)
			utils.ResponseHandler(c, http.StatusInternalServerError, false, "Error: "+msg, nil)
			return
		}

		token, refreshToken, _ := config.JwtWrapper.TokenGenerator(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, *foundUser.User_ID)

		config.JwtWrapper.UpdateAllTokens(token, refreshToken, *foundUser.User_ID)

		utils.ResponseHandler(c, http.StatusFound, true, "Login Successfully !", foundUser)
		ctx.Done()
	}
}

// bson.D{}  -->> D is an ordered representation of a BSON document. This type should be used when the order of the elements matters, such as MongoDB command documents. --->> (Slice)

// bson.M{}  -->>  M is an unordered representation of a BSON document. This type should be used when the order of the elements does not matter. This type is handled as a regular map[string]interface{} when encoding and decoding. --->> (Map)

// bson.A{}: An ordered representation of a BSON array

// bson.E{}: A single element inside a D type
