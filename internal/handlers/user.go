package handlers

import (
	"net/http"

	"github.com/gocli/social_api/internal/services"
	"github.com/gocli/social_api/internal/utils"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	BaseHandler
	userService *services.UserService
	validator   *utils.Validator
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *services.UserService, validator *utils.Validator) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
	}
}

// GetUserProfile handles getting a user's public profile
func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from path
	userID, err := h.ParseIDFromPath(r, "userId")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	// Get user
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	// Return user data
	utils.SendJSONResponse(w, http.StatusOK, user)
}

// SearchUsers handles searching for users
func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query().Get("q")
	limit := h.ParseQueryInt(r, "limit", 20)
	offset := h.ParseQueryInt(r, "offset", 0)

	// Search users
	users, err := h.userService.SearchUsers(query, limit, offset)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Return users
	utils.SendJSONResponse(w, http.StatusOK, users)
}

// GetMe handles getting the current user's profile
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Get user
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	// Return user data
	utils.SendJSONResponse(w, http.StatusOK, user)
}

// UpdateMe handles updating the current user's profile
func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse request body
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := h.DecodeJSONBody(w, r, &req); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	// Get current user
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	// Update user fields
	user.Name = req.Name
	user.Email = req.Email

	// Update user
	if err := h.userService.UpdateUser(user); err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Return updated user
	utils.SendJSONResponse(w, http.StatusOK, user)
}

// PartialUpdateMe handles partially updating the current user's profile
func (h *UserHandler) PartialUpdateMe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse request body
	var req struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
	}

	if err := h.DecodeJSONBody(w, r, &req); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	// Get current user
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	// Update user fields if provided
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	// Update user
	if err := h.userService.UpdateUser(user); err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Return updated user
	utils.SendJSONResponse(w, http.StatusOK, user)
}
