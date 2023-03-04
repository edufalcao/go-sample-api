package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// GetMongoDbConnection get connection of mongodb
func getMongoDbConnection(connString string) (*mongo.Client, error) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connString))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}

func getMongoDbCollection(connString string, dbName string, collectionName string) (*mongo.Collection, error) {
	client, err := getMongoDbConnection(connString)

	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)

	return collection, nil
}
