package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	Id int	`json:"id"`
	Compleled bool `json:"completed"`
	Body string `json:"body"`
}

func main() {
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error load .env file")
	}

	PORT := os.Getenv("PORT")
	
	todos := []Todo{}

// CREATE
	app.Post("/api/todos/", func(c *fiber.Ctx) error {
		todo := &Todo{}

		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(200).JSON(fiber.Map{"message": "To do Body is required"})
		}

		todo.Id = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)
	})

// UPDATE
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.Id) == id {
				todos[i].Compleled = true
				return c.Status(200).JSON(todos)
			}
		}

		return c.Status(404).JSON(fiber.Map{"message": "todos not found"})

	})

// DELETE
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.Id) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"message": "success"})
			}
		}

		return c.Status(404).JSON(fiber.Map{"message": "todos not found"})

	})

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})
	
	log.Fatal(app.Listen(":" + PORT))
}