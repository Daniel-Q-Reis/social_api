package models

// FriendRequest represents a friend request between users
type FriendRequest struct {
	BaseModel
	FromUserID int    `json:"from_user_id" db:"from_user_id"`
	ToUserID   int    `json:"to_user_id" db:"to_user_id"`
	Status     string `json:"status" db:"status"` // pending, accepted, rejected
}

// Friend represents a friendship connection between users
type Friend struct {
	BaseModel
	UserID   int `json:"user_id" db:"user_id"`
	FriendID int `json:"friend_id" db:"friend_id"`
}
