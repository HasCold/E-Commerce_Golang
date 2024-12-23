package utils

import (
	"ecommerce/helpers"
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
		log.Println(err)
		// Triggers a panic: After logging the error, it calls the panic() function, which stops the normal execution of the program and begins the unwinding of the stack. This allows deferred functions to execute before the program terminates.
		log.Panic(err)

		c.JSON(statusCode, gin.H{
			"success": success,
			"message": message,
		})

		return
	}
}

func ResponseHandler(c *gin.Context, statusCode int, success bool, message string, body interface{}) {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	response := gin.H{
		"success": success,
	}

	if message != "" {
		response["message"] = message
	}

	// Check if body is not nil or is not empty
	isEmpty := helpers.IsEmpty(body)
	if body != nil && !isEmpty {
		response["data"] = body
	}

	if c != nil {
		c.JSON(statusCode, response)
	}

	if success == false {
		return
	}
}

// Data abstraction means exposing only the necessary details of an object or module while hiding its implementation. This simplifies the interaction with the object or module and keeps the internal logic hidden.

// Encapsulation										|		Abstraction
// Hides data by restricting access.							Hides implementation details.
// Achieved with exported/unexported fields and methods.		Achieved using interfaces and structs.
// Protects the state of an object.								Focuses on what an object does, not how.
