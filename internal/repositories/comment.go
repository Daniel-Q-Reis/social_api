// Package repositories provides data access functionality for the social media API
package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gocli/social_api/internal/models"
)

// CommentRepository provides methods for accessing comment data
type CommentRepository struct {
	*BaseRepository
}

// NewCommentRepository creates a new CommentRepository
func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{BaseRepository: NewBaseRepository(db)}
}

// CreateComment inserts a new comment into the database
func (r *CommentRepository) CreateComment(comment *models.Comment) error {
	query := `
		INSERT INTO comments (user_id, resource_type, resource_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(query, comment.UserID, comment.ResourceType, comment.ResourceID,
		comment.Content, now, now).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

// GetCommentByID retrieves a comment by ID
func (r *CommentRepository) GetCommentByID(id int) (*models.Comment, error) {
	comment := &models.Comment{}
	query := `
		SELECT id, user_id, resource_type, resource_id, content, created_at, updated_at
		FROM comments
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(&comment.ID, &comment.UserID, &comment.ResourceType,
		&comment.ResourceID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	return comment, nil
}

// GetCommentsForResource retrieves all comments for a specific resource
func (r *CommentRepository) GetCommentsForResource(resourceType string, resourceID int) ([]*models.Comment, error) {
	comments := []*models.Comment{}
	query := `
		SELECT c.id, c.user_id, c.resource_type, c.resource_id, c.content, c.created_at, c.updated_at,
		       u.id, u.name, u.profile_picture_url, u.cover_photo_url, u.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.resource_type = $1 AND c.resource_id = $2
		ORDER BY c.created_at DESC`

	rows, err := r.db.Query(query, resourceType, resourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// In a real application, you would log this error
			// For now, we'll just ignore it to keep the code simple
			_ = closeErr
		}
	}()

	for rows.Next() {
		comment := &models.Comment{}
		user := &models.UserPublic{}
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.ResourceType, &comment.ResourceID,
			&comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
			&user.ID, &user.Name, &user.ProfilePictureURL, &user.CoverPhotoURL, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comment.User = user
		comments = append(comments, comment)
	}

	return comments, nil
}

// UpdateComment updates a comment
func (r *CommentRepository) UpdateComment(comment *models.Comment) error {
	query := `
		UPDATE comments
		SET content = $1, updated_at = $2
		WHERE id = $3`

	now := time.Now()
	_, err := r.db.Exec(query, comment.Content, now, comment.ID)

	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	comment.UpdatedAt = now
	return nil
}

// DeleteComment removes a comment from the database
func (r *CommentRepository) DeleteComment(id int) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}

// GetUserComments retrieves all comments made by a user
func (r *CommentRepository) GetUserComments(userID int) ([]*models.Comment, error) {
	comments := []*models.Comment{}
	query := `
		SELECT id, user_id, resource_type, resource_id, content, created_at, updated_at
		FROM comments
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user comments: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// In a real application, you would log this error
			// For now, we'll just ignore it to keep the code simple
			_ = closeErr
		}
	}()

	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.ResourceType, &comment.ResourceID,
			&comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
