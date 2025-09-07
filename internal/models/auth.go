package models

import (
	"time"
)

// RefreshToken represents a refresh token for a user
type RefreshToken struct {
	BaseModel
	UserID    int       `json:"user_id" db:"user_id"`
	Token     string    `json:"-" db:"token"` // Never serialize token to JSON
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Revoked   bool      `json:"-" db:"revoked"` // Never serialize revoked status to JSON
}
