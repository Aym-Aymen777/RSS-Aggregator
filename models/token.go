package models

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AccessTokenClaims struct {
	UserID bson.ObjectID `json:"_id"`
	Email  string        `json:"email"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID bson.ObjectID `json:"_id"`
	Email  string        `json:"email"`
	jwt.RegisteredClaims
}

// Response struct sent to the client after successful login
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
