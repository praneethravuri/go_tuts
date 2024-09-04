// File: internal/handlers/todo_handlers.go

package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/praneethravuri/go_tuts/react-todo/internal/db"
	"github.com/praneethravuri/go_tuts/react-todo/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllTodos(c *fiber.Ctx) error {
	collection := db.GetCollection("your_database", "todos")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch todos",
		})
	}
	defer cursor.Close(context.Background())

	var todos []models.Todo
	if err := cursor.All(context.Background(), &todos); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode todos",
		})
	}

	return c.JSON(todos)
}