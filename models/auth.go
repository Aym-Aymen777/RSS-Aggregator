package models

import (
	"time"
)


type Auth struct {
Username string `bson:"username" json:"username"`
Email string `bson:"email" json:"email"`
Password string `bson:"password" json:"-"`
CreatedAt time.Time `bson:"created_at" json:"created_at"`
UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}