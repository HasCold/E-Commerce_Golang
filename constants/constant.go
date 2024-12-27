package constants

import (
	"os"
)

var PORT = os.Getenv("PORT")

var SECRET_KEY = os.Getenv("SECRET_KEY")
var ISSUED_BY = os.Getenv("ISSUED_BY")
var EXPIRATION_HOURS = os.Getenv("EXPIRATION_HOURS")
