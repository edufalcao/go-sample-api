package main

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const connString = "mongodb://localhost:27017"
const dbName = "contacts-db"
const collectionName = "contacts"
const port = 8080

func getContact(c *fiber.Ctx) {
	collection, err := getMongoDbCollection(connString, dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
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
		c.Status(500).Send(err)
		return
	}

	cur.All(context.Background(), &results)

	if results == nil {
		c.SendStatus(404)
		return
	}

	json, _ := json.Marshal(results)
	c.Send(json)
}

func createContact(c *fiber.Ctx) {
	collection, err := getMongoDbCollection(connString, dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	var contact Contact
	json.Unmarshal([]byte(c.Body()), &contact)

	res, err := collection.InsertOne(context.Background(), contact)
	if err != nil {
		c.Status(500).Send(err)
		return
	}

	response, _ := json.Marshal(res)
	c.Send(response)
}

func updateContact(c *fiber.Ctx) {
	collection, err := getMongoDbCollection(connString, dbName, collectionName)
	if err != nil {
		c.Status(500).Send(err)
		return
	}
	var contact Contact
	json.Unmarshal([]byte(c.Body()), &contact)

	update := bson.M{
		"$set": contact,
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	response, _ := json.Marshal(res)
	c.Send(response)
}

func deleteContact(c *fiber.Ctx) {
	collection, err := getMongoDbCollection(connString, dbName, collectionName)

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})

	if err != nil {
		c.Status(500).Send(err)
		return
	}

	jsonResponse, _ := json.Marshal(res)
	c.Send(jsonResponse)
}

func main() {
	app := fiber.New()
	app.Get("/contact/:id?", getContact)
	app.Post("/contact", createContact)
	app.Put("/contact/:id", updateContact)
	app.Delete("/contact/:id", deleteContact)
	app.Listen(8080)
}
