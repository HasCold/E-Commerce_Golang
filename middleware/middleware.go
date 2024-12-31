package middleware

import (
	"ecommerce/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		config.TokenSetting()

		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "No authorizaion header provided !"})
			c.Abort()
			return
		}

		if strings.HasPrefix(clientToken, "Bearer") {
			token := strings.Split(clientToken, "Bearer ")

			claims, err := config.JwtWrapper.ValidateToken(strings.TrimSpace(token[1])) // token[0] ==> actual token recived
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
