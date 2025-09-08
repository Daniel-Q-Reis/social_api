package handlers

import (
	"net/http"

	"github.com/gocli/social_api/internal/services"
	"github.com/gocli/social_api/internal/utils"
)

// PostHandler handles post-related HTTP requests
type PostHandler struct {
	BaseHandler
	postService *services.PostService
	validator   *utils.Validator
}

// NewPostHandler creates a new PostHandler
func NewPostHandler(postService *services.PostService, validator *utils.Validator) *PostHandler {
	return &PostHandler{
		postService: postService,
		validator:   validator,
	}
}

// CreatePost handles creating a new post
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse request body
	var req struct {
		Content string `json:"content" validate:"required"`
		Privacy string `json:"privacy"`
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

	// Create post
	post, err := h.postService.CreatePost(userID, req.Content, req.Privacy)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Return post
	utils.SendJSONResponse(w, http.StatusCreated, post)
}

// GetFeed handles getting the user's feed
func (h *PostHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse query parameters
	page := h.ParseQueryInt(r, "page", 1)
	limit := h.ParseQueryInt(r, "limit", 20)
	offset := (page - 1) * limit

	// Get feed
	posts, err := h.postService.GetFeed(userID, limit, offset)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Return posts
	utils.SendJSONResponse(w, http.StatusOK, posts)
}

// GetUserPosts handles getting a user's posts
func (h *PostHandler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from path
	userID, err := h.ParseIDFromPath(r, "userId")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	// Parse query parameters
	page := h.ParseQueryInt(r, "page", 1)
	limit := h.ParseQueryInt(r, "limit", 20)
	offset := (page - 1) * limit

	// Get user's posts
	posts, err := h.postService.GetPostsByUserID(userID, limit, offset)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Return posts
	utils.SendJSONResponse(w, http.StatusOK, posts)
}

// GetPost handles getting a single post
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	// Parse post ID from path
	postID, err := h.ParseIDFromPath(r, "postId")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
		return
	}

	// Get post
	post, err := h.postService.GetPostByID(postID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusNotFound, map[string]string{"error": "Post not found"})
		return
	}

	// Return post
	utils.SendJSONResponse(w, http.StatusOK, post)
}

// UpdatePost handles updating a post
func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse post ID from path
	postID, err := h.ParseIDFromPath(r, "postId")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
		return
	}

	// Parse request body
	var req struct {
		Content string `json:"content"`
		Privacy string `json:"privacy"`
	}

	if err := h.DecodeJSONBody(w, r, &req); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	// Update post
	post, err := h.postService.UpdatePost(postID, userID, req.Content, req.Privacy)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Return post
	utils.SendJSONResponse(w, http.StatusOK, post)
}

// DeletePost handles deleting a post
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse post ID from path
	postID, err := h.ParseIDFromPath(r, "postId")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
		return
	}

	// Delete post
	if err := h.postService.DeletePost(postID, userID); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Return success response
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "Post deleted"})
}
