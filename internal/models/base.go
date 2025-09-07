// Package models provides data structures for the social media API
package models

import (
	"time"
)

// BaseModel represents common fields for all models
type BaseModel struct {
	ID        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
