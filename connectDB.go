package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var MongoClient *mongo.Client

func connectDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dbURI := os.Getenv("MONGODB_URI")
	if dbURI == "" {
		log.Fatal("MONGODB_URI is not set")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().
		ApplyURI(dbURI).
		SetServerAPIOptions(serverAPI)

	// ✅ v2: NO context here
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal(err)
	}

	// ✅ Context is used HERE
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Connected to MongoDB")
	MongoClient = client
	return client
}