package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/praveent04/URL_short/api"
	"github.com/praveent04/URL_short/database"
)

// JWTMiddleware validates JWT tokens
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Printf("JWT Middleware: Missing authorization header for %s %s", c.Method(), c.Path())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			log.Printf("JWT Middleware: Invalid authorization header format for %s %s", c.Method(), c.Path())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header",
			})
		}

		token, err := jwt.ParseWithClaims(tokenString, &api.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil // Use env var
		})

		if err != nil {
			log.Printf("JWT Middleware: Token parse error for %s %s: %v", c.Method(), c.Path(), err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		if !token.Valid {
			log.Printf("JWT Middleware: Token invalid for %s %s", c.Method(), c.Path())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		claims, ok := token.Claims.(*api.JWTClaims)
		if !ok {
			log.Printf("JWT Middleware: Invalid token claims for %s %s", c.Method(), c.Path())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		log.Printf("JWT Middleware: Authenticated user %d for %s %s", claims.UserID, c.Method(), c.Path())
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		return c.Next()
	}
}

func setupRoutes(app *fiber.App) {
	// Public routes (must be first to avoid conflicts)
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "ok",
		})
	})
	app.Get("/api/v1/debug/:url", api.DebugURL) // Debug route to check URLs in Redis

	// Auth routes (public)
	app.Post("/api/v1/register", api.RegisterUser)
	app.Post("/api/v1/login", api.LoginUser)

	// Protected routes
	protected := app.Group("/api/v1", JWTMiddleware())
	protected.Post("/shorten", api.CreateShortURL)
	protected.Get("/urls", api.GetUserURLs)                                // Get user URLs
	protected.Get("/stats/:url", api.GetURLStats)                          // Get URL statistics
	protected.Post("/notifications/send", api.SendExpirationNotifications) // Send expiration notifications

	// Test protected route
	protected.Get("/test", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		userEmail := c.Locals("user_email")
		return c.Status(200).JSON(fiber.Map{
			"message":    "JWT is working!",
			"user_id":    userID,
			"user_email": userEmail,
		})
	})

	// Serve static frontend files (for production deployment)
	app.Static("/", "./frontend/build")

	// Set up redirect route - must be last and should NOT have any middleware
	app.Get("/:url", api.RedirectURL)
}

func main() {
	// Load environment variables
	if _, err := os.Stat(".env"); err == nil {
		// .env file exists, load it
		if loadErr := godotenv.Load(); loadErr != nil {
			log.Printf("Error loading .env file: %v", loadErr)
		} else {
			log.Printf("Successfully loaded .env file")
		}
	} else {
		log.Printf("No .env file found, using system environment variables (normal for production)")
	}

	// Verify essential environment variables
	requiredVars := []string{"DB_ADDR", "DB_HOST_PG", "JWT_SECRET"}
	for _, varName := range requiredVars {
		if os.Getenv(varName) == "" {
			log.Printf("Warning: Required environment variable %s is not set", varName)
		} else {
			log.Printf("Environment variable %s is configured", varName)
		}
	}

	// Initialize Redis
	if err := database.InitRedis(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize PostgreSQL
	if err := database.InitPostgreSQL(); err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
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
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:3001,http://127.0.0.1:3000,http://127.0.0.1:3001,https://ynit.com,http://ynit.com",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

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
