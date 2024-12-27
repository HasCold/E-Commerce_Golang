package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	Order_ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Order_Cart     []ProductUser      `json:"order_list" bson:"order_list"`
	Ordered_At     time.Time          `json:"ordered_at" bson:"ordered_at"`
	Price          int                `json:"total_price" bson:"total_price"`
	Discount       *uint              `json:"discount" bson:"discount"`
	Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool `json:"digital" bson:"digital"`
	COD     bool `json:"cod" bson:"cod"`
}
