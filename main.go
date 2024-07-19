package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	Id        primitive.ObjectID 	`json:"id,omitempty" bson:"_id,omitempty"`
	Completed bool   							`json:"completed"`
	Body      string 							`json:"body"`
}

var collection *mongo.Collection

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("erro to load env", err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOption := options.Client().ApplyURI(MONGODB_URI)
	
	client, err := mongo.Connect(context.Background(), clientOption)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MONGO DB Atlas")

	collection = client.Database("golang_db").Collection("todo")

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodos)
	app.Patch("/api/todos/:id", updateTodos)
	app.Delete("/api/todos/:id", deleteTodos)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + PORT))

}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}

		todos = append(todos, todo)

	}

	return c.JSON(todos)
}

func createTodos(c *fiber.Ctx) error {
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		c.Status(400).JSON(fiber.Map{"message": "Todo body cannot be empty"})
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.Id = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(200).JSON(todo)
}

func updateTodos(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "id not found"})
	}

	filter := bson.M{"_id":objectId}
	update := bson.M{"$set":bson.M{"completed":true}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"message": "update successfully"})
}

func deleteTodos(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "id not found"})
	}

	filter := bson.M{"_id":objectId}

	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil
	}

	return c.Status(200).JSON(fiber.Map{"massage": "delete successfully"})
}