package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// c is a pointer to a gin.Context Struct
// utils.ErrorHandler(err, c, http.StatusInternalServerError, false, err)
func ErrorHandler(err error, c *gin.Context, statusCode int, success bool, message string) {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	if message == "" {
		message = "An unexpected error occurred"
	}

	if err != nil {
		// Triggers a panic: After logging the error, it calls the panic() function, which stops the normal execution of the program and begins the unwinding of the stack. This allows deferred functions to execute before the program terminates.
		log.Panic(err)

		c.JSON(statusCode, gin.H{
			"success": success,
			"message": message,
		})

		return
	}
}

func ResponseHandler(c *gin.Context, statusCode int, success bool, message string) {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	if message == "" {
		message = "An unexpected error occurred"
	}

	if c != nil {
		c.JSON(statusCode, gin.H{
			"success": success,
			"message": message,
		})
	}

	if success == false {
		return
	}
}
