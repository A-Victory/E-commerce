package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

var Client *mongo.Client = dbSetUp()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	coll := client.Database("E-commerce").Collection(collectionName)
	return coll
}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	coll := client.Database("E-commerce").Collection(collectionName)
	return coll
}
