package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Aym-Aymen777/RSS-Aggregator/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var MongoClient *mongo.Client

//insert a user to the collection

func InsertUser(user models.User) error {
	collection := MongoClient.Database("rssagg").Collection("users")
	_, err := collection.InsertOne(context.TODO(), user)
	return err
}

func InsertMany(docs []models.User) error {
	coll := MongoClient.Database("rssagg").Collection("users")
	result, err := coll.InsertMany(context.TODO(), docs)
	if err != nil {
		fmt.Printf("A bulk write error occurred, but %v documents were still inserted.\n", len(result.InsertedIDs))
	}
	for _, id := range result.InsertedIDs {
		fmt.Printf("Inserted document with _id: %v\n", id)
	}
	return err
}

func FindByQuery(query string, value any) []models.User {
	coll := MongoClient.Database("rssagg").Collection("users")
	filter := bson.D{{Key: query, Value: value}}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []models.User
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		res, _ := bson.MarshalExtJSON(result, false, false)
		fmt.Println(string(res))
	}
	return results
}

func UpdateUser(id string) {
	coll := MongoClient.Database("rssagg").Collection("users")

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal("Invalid ObjectID:", err)
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: "Aymen Updated"},
			{Key: "updatedAt", Value: time.Now()},
		}},
		{Key: "$inc", Value: bson.D{
			{Key: "age", Value: 1},
		}},
	}

	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Documents matched: %d\n", result.MatchedCount)
	fmt.Printf("Documents updated: %d\n", result.ModifiedCount)
}
