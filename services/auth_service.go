package services

import (
	"context"
	"log"
	"time"

	"github.com/Aym-Aymen777/RSS-Aggregator/models"
	"github.com/Aym-Aymen777/RSS-Aggregator/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func RegisterUser(col *mongo.Collection, username, email, password string) error {
	ctx := context.Background()
	countMail, _ := col.CountDocuments(ctx, bson.M{"email": email})
	if countMail > 0 {
		log.Println("Email already exists")
		return nil
	}
	countUser, _ := col.CountDocuments(ctx, bson.M{"username": username})
	if countUser > 0 {
		log.Println("Username already exists")
		return nil
	}
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	 newUser := models.Auth{
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = col.InsertOne(ctx, newUser)
	return err
}
