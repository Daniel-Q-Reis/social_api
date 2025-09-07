package models

// Album represents a photo album
type Album struct {
	BaseModel
	UserID      int    `json:"user_id" db:"user_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Privacy     string `json:"privacy" db:"privacy"` // public, friends, only_me
}

// Photo represents a photo in an album
type Photo struct {
	BaseModel
	AlbumID int    `json:"album_id" db:"album_id"`
	URL     string `json:"url" db:"url"`
	Caption string `json:"caption" db:"caption"`
}
