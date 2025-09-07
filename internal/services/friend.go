package services

import (
	"fmt"

	"github.com/gocli/social_api/internal/models"
	"github.com/gocli/social_api/internal/repositories"
)

// FriendService provides friend-related functionality
type FriendService struct {
	BaseService
	friendRepo *repositories.FriendRepository
	userRepo   *repositories.UserRepository
}

// NewFriendService creates a new FriendService
func NewFriendService(friendRepo *repositories.FriendRepository, userRepo *repositories.UserRepository) *FriendService {
	return &FriendService{
		friendRepo: friendRepo,
		userRepo:   userRepo,
	}
}

// SendFriendRequest sends a friend request from one user to another
func (s *FriendService) SendFriendRequest(fromUserID, toUserID int) (*models.FriendRequest, error) {
	// Check if users exist
	_, err := s.userRepo.GetByID(fromUserID)
	if err != nil {
		return nil, fmt.Errorf("sender user not found: %w", err)
	}

	_, err = s.userRepo.GetByID(toUserID)
	if err != nil {
		return nil, fmt.Errorf("recipient user not found: %w", err)
	}

	// Create friend request
	request := &models.FriendRequest{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Status:     "pending",
	}

	if err := s.friendRepo.CreateFriendRequest(request); err != nil {
		return nil, fmt.Errorf("failed to create friend request: %w", err)
	}

	return request, nil
}

// GetFriendRequestsForUser retrieves pending friend requests for a user
func (s *FriendService) GetFriendRequestsForUser(userID int) ([]*models.FriendRequest, error) {
	requests, err := s.friendRepo.GetPendingFriendRequestsForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get friend requests: %w", err)
	}
	return requests, nil
}

// AcceptFriendRequest accepts a friend request
func (s *FriendService) AcceptFriendRequest(requestID, userID int) error {
	// Get the friend request
	request, err := s.friendRepo.GetFriendRequestByID(requestID)
	if err != nil {
		return fmt.Errorf("friend request not found: %w", err)
	}

	// Check if the user is the recipient of the request
	if request.ToUserID != userID {
		return fmt.Errorf("user is not authorized to accept this request")
	}

	// Update request status
	if err := s.friendRepo.UpdateFriendRequestStatus(requestID, "accepted"); err != nil {
		return fmt.Errorf("failed to update friend request: %w", err)
	}

	// Create friendship
	friend := &models.Friend{
		UserID:   request.FromUserID,
		FriendID: request.ToUserID,
	}

	if err := s.friendRepo.CreateFriend(friend); err != nil {
		return fmt.Errorf("failed to create friendship: %w", err)
	}

	return nil
}

// RejectFriendRequest rejects a friend request
func (s *FriendService) RejectFriendRequest(requestID, userID int) error {
	// Get the friend request
	request, err := s.friendRepo.GetFriendRequestByID(requestID)
	if err != nil {
		return fmt.Errorf("friend request not found: %w", err)
	}

	// Check if the user is the recipient of the request
	if request.ToUserID != userID {
		return fmt.Errorf("user is not authorized to reject this request")
	}

	// Update request status
	if err := s.friendRepo.UpdateFriendRequestStatus(requestID, "rejected"); err != nil {
		return fmt.Errorf("failed to update friend request: %w", err)
	}

	return nil
}

// GetFriendsForUser retrieves friends for a user
func (s *FriendService) GetFriendsForUser(userID int) ([]*models.UserPublic, error) {
	friends, err := s.friendRepo.GetFriendsForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get friends: %w", err)
	}
	return friends, nil
}

// UnfriendUser removes a friendship between two users
func (s *FriendService) UnfriendUser(userID, friendID int) error {
	// Check if users exist
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	_, err = s.userRepo.GetByID(friendID)
	if err != nil {
		return fmt.Errorf("friend not found: %w", err)
	}

	// Delete friendship
	if err := s.friendRepo.DeleteFriend(userID, friendID); err != nil {
		return fmt.Errorf("failed to unfriend user: %w", err)
	}

	return nil
}
