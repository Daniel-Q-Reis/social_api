// Package models provides data structures for the social media API
package models

import (
	"time"
)

// Like represents a like on a resource (post, photo, etc.)
type Like struct {
	UserID       int       `json:"user_id" db:"user_id"`
	ResourceType string    `json:"resource_type" db:"resource_type"`
	ResourceID   int       `json:"resource_id" db:"resource_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
