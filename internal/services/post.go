package services

import (
	"fmt"

	"github.com/gocli/social_api/internal/models"
	"github.com/gocli/social_api/internal/repositories"
)

// PostService provides post-related functionality
type PostService struct {
	BaseService
	postRepo *repositories.PostRepository
}

// NewPostService creates a new PostService
func NewPostService(postRepo *repositories.PostRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

// CreatePost creates a new post
func (s *PostService) CreatePost(userID int, content, privacy string) (*models.Post, error) {
	post := &models.Post{
		UserID:  userID,
		Content: content,
		Privacy: privacy,
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

// GetPostByID retrieves a post by ID
func (s *PostService) GetPostByID(id int) (*models.Post, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return post, nil
}

// GetPostsByUserID retrieves posts for a specific user
func (s *PostService) GetPostsByUserID(userID int, limit, offset int) ([]*models.Post, error) {
	posts, err := s.postRepo.GetPostsByUserID(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	return posts, nil
}

// GetFeed retrieves posts for a user's feed
func (s *PostService) GetFeed(userID int, limit, offset int) ([]*models.PostWithUser, error) {
	posts, err := s.postRepo.GetFeed(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed: %w", err)
	}
	return posts, nil
}

// UpdatePost updates a post
func (s *PostService) UpdatePost(postID, userID int, content, privacy string) (*models.Post, error) {
	// First get the post to verify ownership
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	// Check if user is authorized to update this post
	if post.UserID != userID {
		return nil, fmt.Errorf("user is not authorized to update this post")
	}

	// Update the post
	post.Content = content
	post.Privacy = privacy

	if err := s.postRepo.Update(post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return post, nil
}

// DeletePost deletes a post
func (s *PostService) DeletePost(postID, userID int) error {
	// First get the post to verify ownership
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	// Check if user is authorized to delete this post
	if post.UserID != userID {
		return fmt.Errorf("user is not authorized to delete this post")
	}

	// Delete the post
	if err := s.postRepo.Delete(postID); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}
