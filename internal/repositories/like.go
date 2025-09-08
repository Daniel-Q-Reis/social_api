// Package repositories provides data access functionality for the social media API
package repositories

import (
	"database/sql"
	"fmt"

	"github.com/gocli/social_api/internal/models"
)

// LikeRepository provides methods for accessing like data
type LikeRepository struct {
	*BaseRepository
}

// NewLikeRepository creates a new LikeRepository
func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{BaseRepository: NewBaseRepository(db)}
}

// CreateLike inserts a new like into the database
func (r *LikeRepository) CreateLike(like *models.Like) error {
	query := `
		INSERT INTO likes (user_id, resource_type, resource_id, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, resource_type, resource_id) DO NOTHING
		RETURNING created_at`

	var createdAt interface{}
	err := r.db.QueryRow(query, like.UserID, like.ResourceType, like.ResourceID, like.CreatedAt).
		Scan(&createdAt)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to create like: %w", err)
	}

	// If ErrNoRows is returned, it means the like already existed (conflict), which is fine
	return nil
}

// DeleteLike removes a like from the database
func (r *LikeRepository) DeleteLike(userID int, resourceType string, resourceID int) error {
	query := `DELETE FROM likes WHERE user_id = $1 AND resource_type = $2 AND resource_id = $3`
	_, err := r.db.Exec(query, userID, resourceType, resourceID)
	if err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}
	return nil
}

// GetLikesForResource retrieves all likes for a specific resource
func (r *LikeRepository) GetLikesForResource(resourceType string, resourceID int) ([]*models.Like, error) {
	likes := []*models.Like{}
	query := `
		SELECT user_id, resource_type, resource_id, created_at
		FROM likes
		WHERE resource_type = $1 AND resource_id = $2
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, resourceType, resourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get likes: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// In a real application, you would log this error
			// For now, we'll just ignore it to keep the code simple
			_ = closeErr
		}
	}()

	for rows.Next() {
		like := &models.Like{}
		err := rows.Scan(&like.UserID, &like.ResourceType, &like.ResourceID, &like.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan like: %w", err)
		}
		likes = append(likes, like)
	}

	return likes, nil
}

// GetUserLikedResources retrieves all resources liked by a user of a specific type
func (r *LikeRepository) GetUserLikedResources(userID int, resourceType string) ([]int, error) {
	resourceIDs := []int{}
	query := `
		SELECT resource_id
		FROM likes
		WHERE user_id = $1 AND resource_type = $2`

	rows, err := r.db.Query(query, userID, resourceType)
	if err != nil {
		return nil, fmt.Errorf("failed to get user liked resources: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// In a real application, you would log this error
			// For now, we'll just ignore it to keep the code simple
			_ = closeErr
		}
	}()

	for rows.Next() {
		var resourceID int
		err := rows.Scan(&resourceID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan resource ID: %w", err)
		}
		resourceIDs = append(resourceIDs, resourceID)
	}

	return resourceIDs, nil
}
