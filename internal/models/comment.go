// Package models provides data structures for the social media API
package models

// Comment represents a comment on a resource (post, photo, etc.)
type Comment struct {
	BaseModel
	UserID       int         `json:"user_id" db:"user_id"`
	ResourceType string      `json:"resource_type" db:"resource_type"`
	ResourceID   int         `json:"resource_id" db:"resource_id"`
	Content      string      `json:"content" db:"content"`
	User         *UserPublic `json:"user,omitempty"` // Populated when fetching comments
}
