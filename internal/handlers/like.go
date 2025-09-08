// Package handlers provides HTTP handlers for the social media API
package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/gocli/social_api/internal/services"
	"github.com/gocli/social_api/internal/utils"
)

// LikeHandler handles like-related HTTP requests
type LikeHandler struct {
	BaseHandler
	likeService *services.LikeService
	validator   *utils.Validator
}

// NewLikeHandler creates a new LikeHandler
func NewLikeHandler(likeService *services.LikeService, validator *utils.Validator) *LikeHandler {
	return &LikeHandler{
		likeService: likeService,
		validator:   validator,
	}
}

// LikeResource handles liking a resource
func (h *LikeHandler) LikeResource(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse resource type and ID from path
	resourceType := chi.URLParam(r, "resourceType")
	resourceIDStr := chi.URLParam(r, "resourceId")
	resourceID, err := strconv.Atoi(resourceIDStr)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid resource ID"})
		return
	}

	// Like the resource
	if err := h.likeService.LikeResource(userID, resourceType, resourceID); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Return success response
	utils.SendJSONResponse(w, http.StatusCreated, map[string]string{"message": "Resource liked"})
}

// UnlikeResource handles unliking a resource
func (h *LikeHandler) UnlikeResource(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse resource type and ID from path
	resourceType := chi.URLParam(r, "resourceType")
	resourceIDStr := chi.URLParam(r, "resourceId")
	resourceID, err := strconv.Atoi(resourceIDStr)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid resource ID"})
		return
	}

	// Unlike the resource
	if err := h.likeService.UnlikeResource(userID, resourceType, resourceID); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Return success response
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "Resource unliked"})
}

// GetLikesForResource handles getting likes for a resource
func (h *LikeHandler) GetLikesForResource(w http.ResponseWriter, r *http.Request) {
	// Parse resource type and ID from path
	resourceType := chi.URLParam(r, "resourceType")
	resourceIDStr := chi.URLParam(r, "resourceId")
	resourceID, err := strconv.Atoi(resourceIDStr)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid resource ID"})
		return
	}

	// Get likes for the resource
	likes, err := h.likeService.GetLikesForResource(resourceType, resourceID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Return likes
	utils.SendJSONResponse(w, http.StatusOK, likes)
}
