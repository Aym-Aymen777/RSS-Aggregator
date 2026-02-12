package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

func NewJWTConfig() *JWTConfig {
	godotenv.Load()
	return &JWTConfig{
		AccessSecret:  os.Getenv("ACCESS_TOKEN_SECRET"),
		RefreshSecret: os.Getenv("REFRESH_TOKEN_SECRET"),
		AccessExpiry:  time.Minute * 15,   //15 min
		RefreshExpiry: time.Hour * 24 * 7, //7 days
	}
}
