package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// primitive.ObjectID is a type defined in the MongoDB Go driver (go.mongodb.org/mongo-driver/bson/primitive). It is used to represent MongoDB's ObjectId, which is the default unique identifier for documents in a MongoDB collection.

type Product struct {
	Product_ID   primitive.ObjectID `json:"product_id,omitempty" bson:"product_id,omitempty"`
	Product_Name *string            `json:"product_name" bson:"product_name"`
	Price        *uint64            `json:"price" bson:"price"`
	Rating       *uint64            `json:"rating" bson:"rating"`
	Image        *string            `json:"image" bson:"image"`
}

type ProductUser struct {
	Product_ID   primitive.ObjectID `json:"product_id,omitempty" bson:"product_id,omitempty"`
	Product_Name *string            `json:"name" bson:"name"`
	Price        *int               `json:"price" bson:"price"`
	Rating       *uint64            `json:"rating" bson:"rating"`
	Image        *string            `json:"image" bson:"image"`
}
