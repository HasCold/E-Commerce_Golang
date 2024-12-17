package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func DBSetup() *mongo.Client {
	uri := "localhost:27017"

	// Golang in-built package Context has some information that may be required to our mongoDB or functions or handlers or routers
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("%v%v", "mongodb://", uri)))
	if err != nil {
		log.Fatal("Error in connecting mongoDB :- ", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Error while taking response from MongoDB :- ", err)
		return nil
	}

	fmt.Println("MongoDB Connected Successfully !!")
	return client
}

var Client *mongo.Client = DBSetup()

// For User Data Collection
func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	var userCollection *mongo.Collection = client.Database("ECommerce").Collection(collectionName)
	return userCollection
}

// For Product Data Collection
func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	var productCollection *mongo.Client = client.Database("ECommerce").Collection(collectionName)
	return productCollection
}
