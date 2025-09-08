// Package handlers provides HTTP handlers for the social media API
package handlers

import (
	"net/http"

	"github.com/gocli/social_api/internal/services"
	"github.com/gocli/social_api/internal/utils"
)

// CommentHandler handles comment-related HTTP requests
type CommentHandler struct {
	BaseHandler
	commentService *services.CommentService
	validator      *utils.Validator
}

// NewCommentHandler creates a new CommentHandler
func NewCommentHandler(commentService *services.CommentService, validator *utils.Validator) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		validator:      validator,
	}
}

// CreateComment handles creating a new comment
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse resource type and ID from path
	// In a real implementation, you would use a router that supports path parameters
	// For simplicity, we'll use placeholder values
	resourceType := "posts" // Placeholder
	resourceID := 1         // Placeholder

	// Parse request body
	var req struct {
		Content string `json:"content" validate:"required"`
	}

	if err := h.DecodeJSONBody(w, r, &req); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	// Validate request
	if err := h.validator.Validate(req); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]interface{}{"error": "Validation failed", "details": err})
		return
	}

	// Create comment
	comment, err := h.commentService.CreateComment(userID, resourceType, resourceID, req.Content)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Return comment
	utils.SendJSONResponse(w, http.StatusCreated, comment)
}

// GetCommentsForResource handles getting comments for a resource
func (h *CommentHandler) GetCommentsForResource(w http.ResponseWriter, r *http.Request) {
	// Parse resource type and ID from path
	// In a real implementation, you would use a router that supports path parameters
	// For simplicity, we'll use placeholder values
	resourceType := "posts" // Placeholder
	resourceID := 1         // Placeholder

	// Get comments for the resource
	comments, err := h.commentService.GetCommentsForResource(resourceType, resourceID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Return comments
	utils.SendJSONResponse(w, http.StatusOK, comments)
}

// DeleteComment handles deleting a comment
func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse comment ID from path
	// In a real implementation, you would use a router that supports path parameters
	// For simplicity, we'll use placeholder values
	commentID := 1 // Placeholder

	// Delete comment
	if err := h.commentService.DeleteComment(commentID, userID); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Return success response
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "Comment deleted"})
}
