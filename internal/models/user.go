package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	BaseModel
	Name              string    `json:"name" db:"name"`
	Email             string    `json:"email" db:"email"`
	Password          string    `json:"-" db:"password"` // Never serialize password to JSON
	BirthDate         time.Time `json:"birth_date" db:"birth_date"`
	ProfilePictureURL string    `json:"profile_picture_url,omitempty" db:"profile_picture_url"`
	CoverPhotoURL     string    `json:"cover_photo_url,omitempty" db:"cover_photo_url"`
}

// UserPublic represents a user's public profile
type UserPublic struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	ProfilePictureURL string    `json:"profile_picture_url,omitempty"`
	CoverPhotoURL     string    `json:"cover_photo_url,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}
