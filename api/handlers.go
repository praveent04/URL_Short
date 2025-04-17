package api

import (
	"fmt"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/praveent04/URL_short/database"
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
	// Parse request body
	body := new(request)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

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

	// Store URL in Redis
	err := database.StoreURL(id, body.URL, body.Expiry)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Custom short URL already in use",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store URL in database",
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
	
	resp := response{
		URL:         body.URL,
		CustomShort: id,
		Expiry:      body.Expiry,
		XRateLimit:  10,
		XRateReset:  30,
		ShortURL:    fmt.Sprintf("%s/%s", domain, id),
	}

	return c.Status(fiber.StatusOK).JSON(resp)
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