package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// dbSetUp implements and instantiates a new MongoDB Client
func dbSetUp() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		fmt.Println("Error connecting to database")
		return nil
	}

	return client
}

// Client is a client for connecting to the database
var Client *mongo.Client = dbSetUp()

// UserData returns a *mongo.Collections for the user's collection
func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	coll := client.Database("E-commerce").Collection(collectionName)
	return coll
}

// UProductData returns a *mongo.Collections for the product's collection
func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	coll := client.Database("E-commerce").Collection(collectionName)
	return coll
}
