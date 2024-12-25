package utils

import (
	"context"
	"ecommerce/database"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserData *mongo.Collection = database.UserData(database.Client, "Users")

type CustomSignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.RegisteredClaims
}

type jwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

func (j *jwtWrapper) TokenGenerator(email, firstName, lastName, userId string) (signedToken string, signedRefreshToken string, err error) {
	claims := &CustomSignedDetails{
		Email:      email,
		First_Name: firstName,
		Last_Name:  lastName,
		Uid:        userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.ExpirationHours) * time.Hour)), // Token expires in 24 hrs
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    j.Issuer,
		},
	}

	refreshClaims := &CustomSignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			// 168 hours (i.e., 7 days) from the current time.
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(168) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedtoken, err := token.SignedString([]byte(j.SecretKey)) // Pass the secret key
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedrefreshtoken, err := refreshToken.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", "", err
	}

	return signedtoken, signedrefreshtoken, nil
}

func (j *jwtWrapper) ValidateToken(signedToken string) (claim *CustomSignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&CustomSignedDetails{},
		func(token *jwt.Token) (interface{}, error) { // The function passed as a third argument returns the secret key used to sign the token, allowing the parsing function to verify the signature.
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		log.Println(err)
		msg = err.Error()
		return
	}

	// we are just checking the token time in this function
	claims, ok := token.Claims.(*CustomSignedDetails)
	if !ok {
		msg = "The token is invalid"
		return
	}

	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		msg = "Token is already expired"
		return
	}

	return claims, msg
}

func (j *jwtWrapper) UpdateAllTokens(signedtoken string, signedrefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updatedObj bson.D

	updatedObj = append(updatedObj, bson.E{Key: "token", Value: signedtoken})
	updatedObj = append(updatedObj, bson.E{Key: "refresh_token", Value: signedrefreshToken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updatedObj = append(updatedObj, bson.E{Key: "updated_at", Value: updated_at})

	// An upsert performs one of the following actions:
	// 1. Updates documents that match your query filter
	// 2. Inserts a new document if there are no matches to your query filter

	upsert := true
	filter := bson.M{"user_id": userId}
	options := options.Update().SetUpsert(upsert)

	_, err := UserData.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updatedObj}}, options)
	if err != nil {
		log.Println(err)
		return
	}

}
