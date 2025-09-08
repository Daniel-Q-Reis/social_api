// Package services provides business logic functionality for the social media API
package services

import (
	"fmt"
	"time"

	"github.com/gocli/social_api/internal/models"
	"github.com/gocli/social_api/internal/repositories"
)

// LikeService provides like-related functionality
type LikeService struct {
	BaseService
	likeRepo *repositories.LikeRepository
}

// NewLikeService creates a new LikeService
func NewLikeService(likeRepo *repositories.LikeRepository) *LikeService {
	return &LikeService{
		likeRepo: likeRepo,
	}
}

// LikeResource creates a like for a resource
func (s *LikeService) LikeResource(userID int, resourceType string, resourceID int) error {
	like := &models.Like{
		UserID:       userID,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		CreatedAt:    time.Now(),
	}

	if err := s.likeRepo.CreateLike(like); err != nil {
		return fmt.Errorf("failed to like resource: %w", err)
	}

	return nil
}

// UnlikeResource removes a like for a resource
func (s *LikeService) UnlikeResource(userID int, resourceType string, resourceID int) error {
	if err := s.likeRepo.DeleteLike(userID, resourceType, resourceID); err != nil {
		return fmt.Errorf("failed to unlike resource: %w", err)
	}

	return nil
}

// GetLikesForResource retrieves all likes for a specific resource
func (s *LikeService) GetLikesForResource(resourceType string, resourceID int) ([]*models.Like, error) {
	likes, err := s.likeRepo.GetLikesForResource(resourceType, resourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get likes: %w", err)
	}

	return likes, nil
}
