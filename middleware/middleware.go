package middleware

import (
	"ecommerce/constants"
	"ecommerce/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var ExpiryTime, _ = strconv.Atoi(constants.EXPIRATION_HOURS)

var jwtWrapper = utils.JWTWrapper{
	SecretKey:       constants.SECRET_KEY,
	Issuer:          constants.ISSUED_BY,
	ExpirationHours: ExpiryTime,
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
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
