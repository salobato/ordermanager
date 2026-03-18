package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Order struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	OrderNumber string        `bson:"order_number"`
	CustomerID  bson.ObjectID `bson:"customer_id"`
	Total       float64       `bson:"total"`
	Status      string        `bson:"status"`
	PlacedAt    time.Time     `bson:"placed_at"`
	UpdatedAt   time.Time     `bson:"updated_at"`
}

func (Order) CollectionName() string {
	return "orders"
}
