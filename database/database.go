package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-mail/mail/v2"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/praveent04/URL_short/models"
)

var (
	client *redis.Client
	ctx    = context.Background()
	db     *gorm.DB
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

// InitPostgreSQL initializes the connection to PostgreSQL
func InitPostgreSQL() error {
	dbHost := getEnv("DB_HOST_PG", "localhost")
	dbPort := getEnv("DB_PORT_PG", "5432")
	dbUser := getEnv("DB_USER_PG", "postgres")
	dbPassword := getEnv("DB_PASS_PG", "")
	dbName := getEnv("DB_NAME", "url_shortener")

	log.Printf("Connecting to PostgreSQL at %s:%s (database: %s, user: %s)", dbHost, dbPort, dbName, dbUser)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.URL{}, &models.Click{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Printf("Successfully connected to PostgreSQL and migrated schema")
	return nil
}

// GetDB returns the GORM database instance
func GetDB() *gorm.DB {
	return db
}

// StoreURLInDB stores URL metadata in PostgreSQL
func StoreURLInDB(userID uint, originalURL, shortCode string, expiryHours uint) (*models.URL, error) {
	url := models.URL{
		UserID:      userID,
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		ExpiryHours: expiryHours,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(time.Duration(expiryHours) * time.Hour),
	}

	result := db.Create(&url)
	if result.Error != nil {
		return nil, result.Error
	}

	return &url, nil
}

// GetURLByShortCode retrieves URL from PostgreSQL by short code
func GetURLByShortCode(shortCode string) (*models.URL, error) {
	var url models.URL
	result := db.Where("short_code = ?", shortCode).First(&url)
	if result.Error != nil {
		return nil, result.Error
	}
	return &url, nil
}

// StoreClick stores a click event in PostgreSQL
func StoreClick(urlID uint, ip, userAgent, country, city, deviceType, browser, os, referrer string) error {
	click := models.Click{
		URLID:      urlID,
		Timestamp:  time.Now(),
		IPAddress:  ip,
		UserAgent:  userAgent,
		Country:    country,
		City:       city,
		DeviceType: deviceType,
		Browser:    browser,
		OS:         os,
		Referrer:   referrer,
	}

	result := db.Create(&click)
	return result.Error
}

// LocationResponse represents the response from ipapi.co
type LocationResponse struct {
	Country string `json:"country_name"`
	City    string `json:"city"`
}

// GetLocationFromIP fetches location data from IP address
func GetLocationFromIP(ip string) (country, city string) {
	if ip == "" || ip == "127.0.0.1" || ip == "::1" {
		return "", ""
	}

	url := fmt.Sprintf("http://ipapi.co/%s/json/", ip)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching location for IP %s: %v", ip, err)
		return "", ""
	}
	defer resp.Body.Close()

	var location LocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		log.Printf("Error decoding location response for IP %s: %v", ip, err)
		return "", ""
	}

	return location.Country, location.City
}

// SendExpirationNotification sends an email notification for URL expiration
func SendExpirationNotification(userEmail, userName, shortCode, originalURL string, expiresAt time.Time) error {
	smtpHost := getEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPortStr := getEnv("SMTP_PORT", "587")
	smtpUser := getEnv("SMTP_USER", "")
	smtpPass := getEnv("SMTP_PASS", "")
	fromEmail := getEnv("FROM_EMAIL", "")
	fromName := getEnv("FROM_NAME", "URL Shortener")

	smtpPort, _ := strconv.Atoi(smtpPortStr)

	if smtpUser == "" || smtpPass == "" {
		log.Printf("Email configuration not set, skipping notification")
		return nil
	}

	// Create email message
	m := mail.NewMessage()
	m.SetHeader("From", m.FormatAddress(fromEmail, fromName))
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Your shortened URL is about to expire")

	// Email body
	body := fmt.Sprintf(`
Hello %s,

Your shortened URL is about to expire!

Short URL: %s/%s
Original URL: %s
Expires on: %s

Please create a new shortened URL if you need to keep this link active.

Best regards,
URL Shortener Team
	`, userName, getEnv("DOMAIN", "short-it.com"), shortCode, originalURL, expiresAt.Format("2006-01-02 15:04:05"))

	m.SetBody("text/plain", body)

	// Send email
	d := mail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Expiration notification sent to %s for URL %s", userEmail, shortCode)
	return nil
}

// GetExpiringURLs returns URLs that are about to expire within the next 24 hours
func GetExpiringURLs() ([]models.URL, error) {
	var urls []models.URL
	now := time.Now()
	expireBefore := now.Add(24 * time.Hour)

	result := db.Preload("User").Where("expires_at BETWEEN ? AND ?", now, expireBefore).Find(&urls)
	if result.Error != nil {
		return nil, result.Error
	}

	return urls, nil
}

// SendExpirationNotifications sends notifications for all expiring URLs
func SendExpirationNotifications() error {
	expiringURLs, err := GetExpiringURLs()
	if err != nil {
		return fmt.Errorf("failed to get expiring URLs: %w", err)
	}

	for _, url := range expiringURLs {
		if url.User.Email != "" {
			err := SendExpirationNotification(
				url.User.Email,
				url.User.Name,
				url.ShortCode,
				url.OriginalURL,
				url.ExpiresAt,
			)
			if err != nil {
				log.Printf("Failed to send notification for URL %s: %v", url.ShortCode, err)
			}
		}
	}

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
