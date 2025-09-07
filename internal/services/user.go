package services

import (
	"fmt"

	"github.com/gocli/social_api/internal/models"
	"github.com/gocli/social_api/internal/repositories"
)

// UserService provides user-related functionality
type UserService struct {
	BaseService
	userRepo *repositories.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id int) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(user *models.User) error {
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// SearchUsers searches for users by query
func (s *UserService) SearchUsers(query string, limit, offset int) ([]*models.UserPublic, error) {
	users, err := s.userRepo.Search(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	return users, nil
}
