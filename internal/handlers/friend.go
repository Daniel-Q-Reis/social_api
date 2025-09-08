package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/gocli/social_api/internal/services"
	"github.com/gocli/social_api/internal/utils"
)

// FriendHandler handles friend-related HTTP requests
type FriendHandler struct {
	BaseHandler
	friendService *services.FriendService
	validator     *utils.Validator
}

// NewFriendHandler creates a new FriendHandler
func NewFriendHandler(friendService *services.FriendService, validator *utils.Validator) *FriendHandler {
	return &FriendHandler{
		friendService: friendService,
		validator:     validator,
	}
}

// GetUserFriends handles getting a user's friends
func (h *FriendHandler) GetUserFriends(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from path
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	// Get friends
	friends, err := h.friendService.GetFriendsForUser(userID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Return friends
	utils.SendJSONResponse(w, http.StatusOK, friends)
}

// GetMyFriendRequests handles getting the current user's friend requests
func (h *FriendHandler) GetMyFriendRequests(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Get friend requests
	requests, err := h.friendService.GetFriendRequestsForUser(userID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Return friend requests
	utils.SendJSONResponse(w, http.StatusOK, requests)
}

// SendFriendRequest handles sending a friend request
func (h *FriendHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (the sender)
	fromUserID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse target user ID from path
	toUserIDStr := chi.URLParam(r, "userId")
	toUserID, err := strconv.Atoi(toUserIDStr)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	// Send friend request
	request, err := h.friendService.SendFriendRequest(fromUserID, toUserID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Return friend request
	utils.SendJSONResponse(w, http.StatusCreated, request)
}

// AcceptFriendRequest handles accepting a friend request
func (h *FriendHandler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse request ID from path
	requestIDStr := chi.URLParam(r, "requestId")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request ID"})
		return
	}

	// Accept friend request
	if err := h.friendService.AcceptFriendRequest(requestID, userID); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Return success response
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "Friend request accepted"})
}

// RejectFriendRequest handles rejecting a friend request
func (h *FriendHandler) RejectFriendRequest(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse request ID from path
	requestIDStr := chi.URLParam(r, "requestId")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request ID"})
		return
	}

	// Reject friend request
	if err := h.friendService.RejectFriendRequest(requestID, userID); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Return success response
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "Friend request rejected"})
}

// UnfriendUser handles unfriending a user
func (h *FriendHandler) UnfriendUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse friend ID from path
	friendIDStr := chi.URLParam(r, "userId")
	friendID, err := strconv.Atoi(friendIDStr)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	// Unfriend user
	if err := h.friendService.UnfriendUser(userID, friendID); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Return success response
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "User unfriended"})
}
