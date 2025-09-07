package models

// Post represents a post created by a user
type Post struct {
	BaseModel
	UserID  int    `json:"user_id" db:"user_id"`
	Content string `json:"content" db:"content"`
	Privacy string `json:"privacy" db:"privacy"` // public, friends, only_me
}

// PostWithUser represents a post with user information
type PostWithUser struct {
	Post
	UserName          string `json:"user_name"`
	ProfilePictureURL string `json:"profile_picture_url,omitempty"`
}
