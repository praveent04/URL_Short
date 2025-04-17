package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
	ctx    = context.Background()
)

// InitRedis initializes the connection to Redis
func InitRedis() error {
	// Get Redis connection details from environment variables
	redisHost := getEnv("DB_ADDR", "localhost")
	redisPort := getEnv("DB_PORT", "6379")
	redisPassword := getEnv("DB_PASS", "")
	redisUser := getEnv("DB_USER", "")

	// Construct Redis URL
	var redisURL string
	if redisUser != "" && redisPassword != "" {
		redisURL = fmt.Sprintf("redis://%s:%s@%s:%s", redisUser, redisPassword, redisHost, redisPort)
	} else if redisPassword != "" {
		redisURL = fmt.Sprintf("redis://:%s@%s:%s", redisPassword, redisHost, redisPort)
	} else {
		redisURL = fmt.Sprintf("redis://%s:%s", redisHost, redisPort)
	}

	log.Printf("Connecting to Redis at %s", redisHost)

	// Create Redis client
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return fmt.Errorf("failed to parse Redis URL: %w", err)
	}
	
	client = redis.NewClient(opt)

	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	
	log.Printf("Successfully connected to Redis at %s:%s", redisHost, redisPort)
	return nil
}

// GetRedisClient returns the Redis client
func GetRedisClient() *redis.Client {
	return client
}

// StoreURL stores a URL in Redis with an expiration time
func StoreURL(id, url string, exp uint) error {
	// Check if ID already exists
	exists, err := client.Exists(ctx, id).Result()
	if err != nil {
		return fmt.Errorf("redis error checking existence: %w", err)
	}
	if exists > 0 {
		return errors.New("URL custom short already exists")
	}

	// Set expiry time (default to 24 hours if not specified)
	expiry := time.Duration(exp) * time.Hour
	if exp == 0 {
		expiry = 24 * time.Hour
	}

	// Log the storage operation
	log.Printf("Storing URL: ID='%s', URL='%s', Expiry=%v", id, url, expiry)

	// Store URL with expiry
	err = client.Set(ctx, id, url, expiry).Err()
	if err != nil {
		return fmt.Errorf("redis error storing URL: %w", err)
	}

	// Verify the URL was stored correctly
	storedURL, err := client.Get(ctx, id).Result()
	if err != nil {
		return fmt.Errorf("redis error verifying storage: %w", err)
	}
	
	log.Printf("Verified storage: ID='%s', StoredURL='%s'", id, storedURL)
	return nil
}

// GetOriginalURL retrieves the original URL from a shortened ID
func GetOriginalURL(id string) (string, error) {
	url, err := client.Get(ctx, id).Result()
	if err == redis.Nil {
		return "", nil // ID not found
	}
	if err != nil {
		return "", fmt.Errorf("redis error: %w", err)
	}
	return url, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// ParseUint parses a string to uint with a default value
func ParseUint(s string, defaultValue uint) uint {
	value, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return defaultValue
	}
	return uint(value)
}