package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Post struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Link        string `bson:"link" json:"link"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}