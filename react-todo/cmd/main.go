// File: cmd/main.go

package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/praneethravuri/go_tuts/react-todo/internal/db"
	"github.com/praneethravuri/go_tuts/react-todo/internal/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	err = db.ConnectToMongoDB(mongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Client.Disconnect(context.Background())

	app := fiber.New()

	// Set up routes
	app.Get("/todos", handlers.GetAllTodos)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}