package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Address struct {
	Address_ID primitive.ObjectID `json:"address_id,omitempty" bson:"address_id,omitempty"`
	House      *string            `json:"house" bson:"house"`
	Street     *string            `json:"street" bson:"street"`
	City       *string            `json:"city" bson:"city"`
	Pincode    *string            `json:"pincode" bson:"pincode"`
}
