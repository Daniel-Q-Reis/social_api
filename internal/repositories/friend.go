package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gocli/social_api/internal/models"
)

// FriendRepository provides methods for accessing friend data
type FriendRepository struct {
	*BaseRepository
}

// NewFriendRepository creates a new FriendRepository
func NewFriendRepository(db *sql.DB) *FriendRepository {
	return &FriendRepository{BaseRepository: NewBaseRepository(db)}
}

// CreateFriendRequest creates a new friend request
func (r *FriendRepository) CreateFriendRequest(request *models.FriendRequest) error {
	query := `
		INSERT INTO friend_requests (from_user_id, to_user_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(query, request.FromUserID, request.ToUserID, request.Status, now, now).
		Scan(&request.ID, &request.CreatedAt, &request.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create friend request: %w", err)
	}

	return nil
}

// GetFriendRequestByID retrieves a friend request by ID
func (r *FriendRepository) GetFriendRequestByID(id int) (*models.FriendRequest, error) {
	request := &models.FriendRequest{}
	query := `
		SELECT id, from_user_id, to_user_id, status, created_at, updated_at
		FROM friend_requests
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(&request.ID, &request.FromUserID, &request.ToUserID,
		&request.Status, &request.CreatedAt, &request.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("friend request not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get friend request: %w", err)
	}

	return request, nil
}

// GetPendingFriendRequestsForUser retrieves pending friend requests for a user
func (r *FriendRepository) GetPendingFriendRequestsForUser(userID int) ([]*models.FriendRequest, error) {
	requests := []*models.FriendRequest{}
	query := `
		SELECT id, from_user_id, to_user_id, status, created_at, updated_at
		FROM friend_requests
		WHERE to_user_id = $1 AND status = 'pending'
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get friend requests: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// In a real application, you would log this error
			// For now, we'll just ignore it to keep the code simple
			_ = closeErr
		}
	}()

	for rows.Next() {
		request := &models.FriendRequest{}
		err := rows.Scan(&request.ID, &request.FromUserID, &request.ToUserID,
			&request.Status, &request.CreatedAt, &request.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan friend request: %w", err)
		}
		requests = append(requests, request)
	}

	return requests, nil
}

// UpdateFriendRequestStatus updates the status of a friend request
func (r *FriendRepository) UpdateFriendRequestStatus(requestID int, status string) error {
	query := `
		UPDATE friend_requests
		SET status = $1, updated_at = $2
		WHERE id = $3`

	now := time.Now()
	_, err := r.db.Exec(query, status, now, requestID)

	if err != nil {
		return fmt.Errorf("failed to update friend request: %w", err)
	}

	return nil
}

// CreateFriend creates a new friendship
func (r *FriendRepository) CreateFriend(friend *models.Friend) error {
	// Create friendship in both directions
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			// In a real application, you would log this error
			// For now, we'll just ignore it to keep the code simple
			_ = rollbackErr
		}
	}()

	query := `
		INSERT INTO friends (user_id, friend_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err = tx.QueryRow(query, friend.UserID, friend.FriendID, now, now).
		Scan(&friend.ID, &friend.CreatedAt, &friend.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create friendship: %w", err)
	}

	// Create reverse friendship
	_, err = tx.Exec(query, friend.FriendID, friend.UserID, now, now)
	if err != nil {
		return fmt.Errorf("failed to create reverse friendship: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetFriendsForUser retrieves friends for a user
func (r *FriendRepository) GetFriendsForUser(userID int) ([]*models.UserPublic, error) {
	friends := []*models.UserPublic{}
	query := `
		SELECT u.id, u.name, u.profile_picture_url, u.cover_photo_url, u.created_at
		FROM friends f
		JOIN users u ON f.friend_id = u.id
		WHERE f.user_id = $1
		ORDER BY u.name`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get friends: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// In a real application, you would log this error
			// For now, we'll just ignore it to keep the code simple
			_ = closeErr
		}
	}()

	for rows.Next() {
		friend := &models.UserPublic{}
		err := rows.Scan(&friend.ID, &friend.Name, &friend.ProfilePictureURL, &friend.CoverPhotoURL, &friend.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan friend: %w", err)
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

// DeleteFriend deletes a friendship
func (r *FriendRepository) DeleteFriend(userID, friendID int) error {
	// Delete friendship in both directions
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			// In a real application, you would log this error
			// For now, we'll just ignore it to keep the code simple
			_ = rollbackErr
		}
	}()

	query := `DELETE FROM friends WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)`
	_, err = tx.Exec(query, userID, friendID)
	if err != nil {
		return fmt.Errorf("failed to delete friendship: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
