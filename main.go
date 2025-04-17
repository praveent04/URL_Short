package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/praveent04/URL_short/api"
	"github.com/praveent04/URL_short/database"
)

func setupRoutes(app *fiber.App) {
	// Set up API routes
	app.Post("/api/v1/shorten", api.CreateShortURL)
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "ok",
		})
	})
	app.Get("/api/v1/debug/:url", api.DebugURL) // Debug route to check URLs in Redis
	
	// Set up redirect route - must be last
	app.Get("/:url", api.RedirectURL)
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		// Continue anyway, we might be using environment variables directly
	}

	// Initialize Redis
	if err := database.InitRedis(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Handle panics and errors
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Setup routes
	setupRoutes(app)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}