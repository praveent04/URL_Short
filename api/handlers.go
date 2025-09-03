package api

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/mssola/user_agent"
	"golang.org/x/crypto/bcrypt"
	"github.com/praveent04/URL_short/database"
	"github.com/praveent04/URL_short/models"
)

// Request model for shortening URLs
type request struct {
	URL         string `json:"url"`
	CustomShort string `json:"custom_short"`
	Expiry      uint   `json:"expiry"`
}

// Response model for shortened URLs
type response struct {
	URL         string `json:"url"`
	ShortURL    string `json:"short_url"`
	CustomShort string `json:"custom_short"`
	Expiry      uint   `json:"expiry"`
	XRateLimit  int    `json:"rate_limit"`
	XRateReset  int    `json:"rate_reset"`
}

// RedirectURL handles redirecting short URLs to their original destination
func RedirectURL(c *fiber.Ctx) error {
    // Get the short URL ID from the URL parameter
    shortID := c.Params("url")

    // Add logging
    fmt.Printf("Received redirect request for ID: '%s'\n", shortID)

    // Get the original URL from Redis
    original, err := database.GetOriginalURL(shortID)
    if err != nil {
        fmt.Printf("Error retrieving URL for %s: %v\n", shortID, err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to retrieve URL",
            "details": err.Error(),
        })
    }

    // If original URL is not found, return a 404
    if original == "" {
        fmt.Printf("URL not found for ID: '%s'\n", shortID)
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "URL not found",
        })
    }

    // Get URL metadata from PostgreSQL for analytics
    urlModel, err := database.GetURLByShortCode(shortID)
    if err != nil {
        fmt.Printf("Error retrieving URL metadata for %s: %v\n", shortID, err)
        // Continue with redirect even if metadata fetch fails
    } else {
        // Parse user agent
        ua := user_agent.New(c.Get("User-Agent"))
        browser, version := ua.Browser()

        // Get IP address
        ip := c.IP()
        if ip == "" {
            ip = c.Get("X-Forwarded-For")
        }
        if ip == "" {
            ip = c.Get("X-Real-IP")
        }

        // Get location from IP
        country, city := database.GetLocationFromIP(ip)

        // Determine device type
        deviceType := "desktop"
        if ua.Mobile() {
            deviceType = "mobile"
        } else if ua.Bot() {
            deviceType = "bot"
        }

        os := ua.OS()

        referrer := c.Get("Referer")

        // Store click in database
        err = database.StoreClick(urlModel.ID, ip, c.Get("User-Agent"), country, city, deviceType, browser+" "+version, os, referrer)
        if err != nil {
            fmt.Printf("Error storing click for %s: %v\n", shortID, err)
            // Don't fail the redirect
        }
    }

    fmt.Printf("Successfully found URL. Redirecting '%s' to '%s'\n", shortID, original)

    // Ensure the URL has a protocol
    if !strings.HasPrefix(original, "http://") && !strings.HasPrefix(original, "https://") {
        original = "http://" + original
    }

    // Use 302 Found for redirects and set explicit Location header
    c.Set("Location", original)
    return c.SendStatus(fiber.StatusFound) // 302 Found
}

// CreateShortURL handles shortening of URLs
func CreateShortURL(c *fiber.Ctx) error {
	log.Printf("CreateShortURL: Received request from user %v", c.Locals("user_id"))

	// Parse request body
	body := new(request)
	if err := c.BodyParser(body); err != nil {
		log.Printf("CreateShortURL: Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	log.Printf("CreateShortURL: Parsed request - URL: %s, CustomShort: %s, Expiry: %d", body.URL, body.CustomShort, body.Expiry)

	// Validate URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	// Add http:// prefix if not present
	if !strings.HasPrefix(body.URL, "http://") && !strings.HasPrefix(body.URL, "https://") {
		body.URL = "http://" + body.URL
	}

	// Generate short ID if not provided
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	// Validate expiry (default to 24 hours)
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	// Get user ID from context
	userID := c.Locals("user_id").(uint)

	// Store URL in PostgreSQL
	urlModel, err := database.StoreURLInDB(userID, body.URL, id, body.Expiry)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Custom short URL already in use",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store URL in database",
			"details": err.Error(),
		})
	}

	// Store URL in Redis for fast lookup
	err = database.StoreURL(id, body.URL, body.Expiry)
	if err != nil {
		// If Redis fails, delete from PG and return error
		database.GetDB().Delete(urlModel)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store URL in cache",
			"details": err.Error(),
		})
	}

	// Prepare response
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = c.Hostname()
	}

	// Make sure domain has protocol
	if !strings.HasPrefix(domain, "http") {
		domain = "http://" + domain
	}

	// Return response with all the fields the frontend expects
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":           urlModel.ID,
		"short_code":   urlModel.ShortCode,
		"original_url": urlModel.OriginalURL,
		"short_url":    fmt.Sprintf("%s/%s", domain, urlModel.ShortCode),
		"expiry":       urlModel.ExpiryHours,
		"created_at":   urlModel.CreatedAt,
		"expires_at":   urlModel.ExpiresAt,
		"url":          body.URL, // Keep for backward compatibility
		"custom_short": id,       // Keep for backward compatibility
		"rate_limit":   10,
		"rate_reset":   30,
	})
}

// DebugURL provides a way to check what's stored in Redis for a given ID
func DebugURL(c *fiber.Ctx) error {
    shortID := c.Params("url")
    
    // Get the original URL from Redis
    original, err := database.GetOriginalURL(shortID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to retrieve URL",
            "details": err.Error(),
        })
    }

    if original == "" {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "URL not found for ID: " + shortID,
        })
    }

    // Return the original URL as JSON for debugging
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "id": shortID,
        "original_url": original,
    })
}

// GetURLStats returns statistics for a specific URL
func GetURLStats(c *fiber.Ctx) error {
    shortID := c.Params("url")
    log.Printf("GetURLStats: Request for short code: %s", shortID)

    // Get URL from database
    urlModel, err := database.GetURLByShortCode(shortID)
    if err != nil {
        log.Printf("GetURLStats: URL not found for short code: %s, error: %v", shortID, err)
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "URL not found",
        })
    }

    log.Printf("GetURLStats: Found URL with ID: %d", urlModel.ID)

    // Get click count
    var clickCount int64
    database.GetDB().Model(&models.Click{}).Where("url_id = ?", urlModel.ID).Count(&clickCount)

    // Get clicks by date (last 30 days)
    var clicksByDate []struct {
        Date  string `json:"date"`
        Count int64  `json:"count"`
    }
    database.GetDB().Model(&models.Click{}).
        Select("DATE(timestamp) as date, COUNT(*) as count").
        Where("url_id = ? AND timestamp >= ?", urlModel.ID, time.Now().AddDate(0, 0, -30)).
        Group("DATE(timestamp)").
        Order("date DESC").
        Scan(&clicksByDate)

    // Get top countries
    var topCountries []struct {
        Country string `json:"country"`
        Count   int64  `json:"count"`
    }
    database.GetDB().Model(&models.Click{}).
        Select("country, COUNT(*) as count").
        Where("url_id = ? AND country != ''", urlModel.ID).
        Group("country").
        Order("count DESC").
        Limit(10).
        Scan(&topCountries)

    // Get device types
    var deviceStats []struct {
        DeviceType string `json:"device_type"`
        Count      int64  `json:"count"`
    }
    database.GetDB().Model(&models.Click{}).
        Select("device_type, COUNT(*) as count").
        Where("url_id = ?", urlModel.ID).
        Group("device_type").
        Scan(&deviceStats)

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "url": fiber.Map{
            "id":          urlModel.ID,
            "short_code":  urlModel.ShortCode,
            "original_url": urlModel.OriginalURL,
            "created_at":  urlModel.CreatedAt,
            "expires_at":  urlModel.ExpiresAt,
        },
        "stats": fiber.Map{
            "total_clicks":   clickCount,
            "clicks_by_date": clicksByDate,
            "top_countries":  topCountries,
            "device_stats":   deviceStats,
        },
    })
}

// GetUserURLs returns all URLs for a user
func GetUserURLs(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)

    var urls []models.URL
    database.GetDB().Where("user_id = ?", userID).Find(&urls)

    // Convert to frontend-friendly format
    var formattedUrls []fiber.Map
    domain := os.Getenv("DOMAIN")
    if domain == "" {
        domain = "http://localhost:3000"
    }

    for _, url := range urls {
        formattedUrls = append(formattedUrls, fiber.Map{
            "id":           url.ID,
            "short_code":   url.ShortCode,
            "original_url": url.OriginalURL,
            "short_url":    fmt.Sprintf("%s/%s", domain, url.ShortCode),
            "expiry":       url.ExpiryHours,
            "created_at":   url.CreatedAt,
            "expires_at":   url.ExpiresAt,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "urls": formattedUrls,
    })
}

// JWTClaims represents JWT claims
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterUser handles user registration
func RegisterUser(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create user
	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	result := database.GetDB().Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User already exists",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// LoginUser handles user login
func LoginUser(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	// Find user
	var user models.User
	result := database.GetDB().Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Generate JWT
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your-secret-key")) // Use env var in production
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": tokenString,
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// SendExpirationNotifications triggers sending of expiration notifications
func SendExpirationNotifications(c *fiber.Ctx) error {
	// This could be run as a background job in production
	err := database.SendExpirationNotifications()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send expiration notifications",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Expiration notifications sent successfully",
	})
}