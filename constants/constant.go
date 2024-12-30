package constants

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	PORT             string
	MONGO_URI        string
	SECRET_KEY       string
	ISSUED_BY        string
	EXPIRATION_HOURS string
)

// Initialize the environment variables once
func LoadENV() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file, falling back to defaults")
	} else {
		log.Println("Environment variables loaded successfully")
	}

	PORT = os.Getenv("PORT")
	MONGO_URI = os.Getenv("MONGO_URI")
	SECRET_KEY = os.Getenv("SECRET_KEY")
	ISSUED_BY = os.Getenv("ISSUED_BY")
	EXPIRATION_HOURS = os.Getenv("EXPIRATION_HOURS")
}
