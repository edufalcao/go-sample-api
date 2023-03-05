package main

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Move to config
const connString = "mongodb://localhost:27017"
const dbName = "contacts-db"
const collectionName = "contacts"
const port = ":8080"

func getContact(c *fiber.Ctx) error {
	collection, err := getMongoDbCollection(connString, dbName, collectionName)

	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	var filter bson.M = bson.M{}

	if c.Params("id") != "" {
		id := c.Params("id")
		objID, _ := primitive.ObjectIDFromHex(id)
		filter = bson.M{"_id": objID}
	}

	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())

	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	cur.All(context.Background(), &results)

	if results == nil {
		return fiber.NewError(404)
	}

	response, _ := json.Marshal(results)
	return c.SendString(string(response))
}

func createContact(c *fiber.Ctx) error {
	collection, err := getMongoDbCollection(connString, dbName, collectionName)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	var contact Contact
	json.Unmarshal([]byte(c.Body()), &contact)

	res, err := collection.InsertOne(context.Background(), contact)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	response, _ := json.Marshal(res)
	return c.SendString(string(response))
}

func updateContact(c *fiber.Ctx) error {
	collection, err := getMongoDbCollection(connString, dbName, collectionName)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}
	var contact Contact
	json.Unmarshal([]byte(c.Body()), &contact)

	update := bson.M{
		"$set": contact,
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)

	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	response, _ := json.Marshal(res)
	return c.SendString(string(response))
}

func deleteContact(c *fiber.Ctx) error {
	collection, err := getMongoDbCollection(connString, dbName, collectionName)

	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})

	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	response, _ := json.Marshal(res)
	return c.SendString(string(response))
}

func main() {
	app := fiber.New()

	app.Get("/contact/:id?", getContact)
	app.Post("/contact", createContact)
	app.Put("/contact/:id", updateContact)
	app.Delete("/contact/:id", deleteContact)

	app.Listen(port)
}
