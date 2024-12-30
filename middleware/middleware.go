package middleware

import (
	"ecommerce/constants"
	"ecommerce/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var SECRET_KEY string
var ISSUED_BY string
var EXPIRATION_HOURS string
var ExpiryTime int

func config() {
	constants.LoadENV()
	SECRET_KEY = constants.SECRET_KEY
	ISSUED_BY = constants.ISSUED_BY
	EXPIRATION_HOURS = constants.EXPIRATION_HOURS
	ExpiryTime, _ = strconv.Atoi(EXPIRATION_HOURS)
}

var jwtWrapper = utils.JWTWrapper{
	SecretKey:       SECRET_KEY,
	Issuer:          ISSUED_BY,
	ExpirationHours: ExpiryTime,
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		config()

		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "No authorizaion header provided !"})
			c.Abort()
			return
		}

		if strings.HasPrefix(clientToken, "Bearer") {
			token := strings.Split(clientToken, "Bearer ")

			claims, err := jwtWrapper.ValidateToken(token[0]) // token[0] ==> actual token recived
			if err != "" {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				c.Abort()
				return
			}

			// Adding the things into the request object via middleware operations
			c.Set("email", claims.Email)
			c.Set("uid", claims.Uid)
			c.Next()

		} else {
			c.JSON(401, gin.H{"error": "Not Authorized !"})
			c.Abort()
			return
		}
	}
}
