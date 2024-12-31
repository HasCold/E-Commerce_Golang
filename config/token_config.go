package config

import (
	"ecommerce/constants"
	"ecommerce/utils"
	"strconv"
)

var JwtWrapper utils.JWTWrapper

func TokenSetting() {
	constants.LoadENV()

	var expiryTime int
	expiryTime, _ = strconv.Atoi(constants.EXPIRATION_HOURS)

	JwtWrapper = utils.JWTWrapper{
		SecretKey:       constants.SECRET_KEY,
		Issuer:          constants.ISSUED_BY,
		ExpirationHours: expiryTime,
	}
}
