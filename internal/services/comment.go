// Package services provides business logic functionality for the social media API
package services

import (
	"fmt"

	"github.com/gocli/social_api/internal/models"
	"github.com/gocli/social_api/internal/repositories"
)

// CommentService provides comment-related functionality
type CommentService struct {
	BaseService
	commentRepo *repositories.CommentRepository
	userRepo    *repositories.UserRepository
}

// NewCommentService creates a new CommentService
func NewCommentService(commentRepo *repositories.CommentRepository, userRepo *repositories.UserRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		userRepo:    userRepo,
	}
}

// CreateComment creates a new comment for a resource
func (s *CommentService) CreateComment(userID int, resourceType string, resourceID int, content string) (*models.Comment, error) {
	comment := &models.Comment{
		UserID:       userID,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Content:      content,
	}

	if err := s.commentRepo.CreateComment(comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Get user info to populate in response
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	comment.User = &models.UserPublic{
		ID:                user.ID,
		Name:              user.Name,
		ProfilePictureURL: user.ProfilePictureURL,
		CoverPhotoURL:     user.CoverPhotoURL,
		CreatedAt:         user.CreatedAt,
	}

	return comment, nil
}

// GetCommentsForResource retrieves all comments for a specific resource
func (s *CommentService) GetCommentsForResource(resourceType string, resourceID int) ([]*models.Comment, error) {
	comments, err := s.commentRepo.GetCommentsForResource(resourceType, resourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	return comments, nil
}

// DeleteComment deletes a comment
func (s *CommentService) DeleteComment(commentID, userID int) error {
	// First get the comment to verify ownership
	comment, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		return fmt.Errorf("failed to get comment: %w", err)
	}

	// Check if user is authorized to delete this comment
	// Either the comment author or the resource owner can delete
	if comment.UserID != userID {
		return fmt.Errorf("user is not authorized to delete this comment")
	}

	// Delete the comment
	if err := s.commentRepo.DeleteComment(commentID); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}
