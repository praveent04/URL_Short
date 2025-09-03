package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user account
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"` // Don't serialize password
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	URLs      []URL          `json:"urls" gorm:"foreignKey:UserID"`
}

// URL represents a shortened URL
type URL struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id"`
	OriginalURL string    `json:"original_url" gorm:"not null"`
	ShortCode   string    `json:"short_code" gorm:"uniqueIndex;not null"`
	ExpiryHours uint      `json:"expiry_hours"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Clicks      []Click   `json:"clicks" gorm:"foreignKey:URLID"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
}

// Click represents a click on a shortened URL for analytics
type Click struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	URLID       uint      `json:"url_id"`
	Timestamp   time.Time `json:"timestamp"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Country     string    `json:"country"`
	City        string    `json:"city"`
	DeviceType  string    `json:"device_type"`
	Browser     string    `json:"browser"`
	OS          string    `json:"os"`
	Referrer    string    `json:"referrer"`
	URL         URL       `json:"url" gorm:"foreignKey:URLID"`
}

// TableName overrides the table name for Click
func (Click) TableName() string {
	return "clicks"
}