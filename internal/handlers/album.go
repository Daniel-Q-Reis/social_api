package handlers

import (
	"net/http"

	"github.com/gocli/social_api/internal/services"
	"github.com/gocli/social_api/internal/utils"
)

// AlbumHandler handles album-related HTTP requests
type AlbumHandler struct {
	BaseHandler
	albumService *services.AlbumService
	validator    *utils.Validator
}

// NewAlbumHandler creates a new AlbumHandler
func NewAlbumHandler(albumService *services.AlbumService, validator *utils.Validator) *AlbumHandler {
	return &AlbumHandler{
		albumService: albumService,
		validator:    validator,
	}
}

// CreateAlbum handles creating a new album
func (h *AlbumHandler) CreateAlbum(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse request body
	var req struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		Privacy     string `json:"privacy"`
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

	// Create album
	album, err := h.albumService.CreateAlbum(userID, req.Name, req.Description, req.Privacy)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Return album
	utils.SendJSONResponse(w, http.StatusCreated, album)
}

// GetUserAlbums handles getting a user's albums
func (h *AlbumHandler) GetUserAlbums(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from path
	userID, err := h.ParseIDFromPath(r, "userId")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	// Get user's albums
	albums, _, err := h.albumService.GetAlbumsWithPhotosByUserID(userID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Return albums
	utils.SendJSONResponse(w, http.StatusOK, albums)
}

// GetAlbum handles getting an album
func (h *AlbumHandler) GetAlbum(w http.ResponseWriter, r *http.Request) {
	// Parse album ID from path
	albumID, err := h.ParseIDFromPath(r, "albumId")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid album ID"})
		return
	}

	// Get album
	album, err := h.albumService.GetAlbumByID(albumID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusNotFound, map[string]string{"error": "Album not found"})
		return
	}

	// Return album
	utils.SendJSONResponse(w, http.StatusOK, album)
}

// UpdateAlbum handles updating an album
func (h *AlbumHandler) UpdateAlbum(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse album ID from path
	albumID, err := h.ParseIDFromPath(r, "albumId")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid album ID"})
		return
	}

	// Parse request body
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Privacy     string `json:"privacy"`
	}

	if err := h.DecodeJSONBody(w, r, &req); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	// Update album
	album, err := h.albumService.UpdateAlbum(albumID, userID, req.Name, req.Description, req.Privacy)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Return album
	utils.SendJSONResponse(w, http.StatusOK, album)
}

// DeleteAlbum handles deleting an album
func (h *AlbumHandler) DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse album ID from path
	albumID, err := h.ParseIDFromPath(r, "albumId")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid album ID"})
		return
	}

	// Delete album
	if err := h.albumService.DeleteAlbum(albumID, userID); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Return success response
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "Album deleted"})
}
