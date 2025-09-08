package handlers

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

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
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.Atoi(userIDStr)
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

// UploadProfilePicture handles uploading a profile picture
func (h *UserHandler) UploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse multipart form with max memory of 10MB
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Unable to parse form"})
		return
	}

	// Get the file from the form
	file, handler, err := r.FormFile("profile_picture")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Unable to get profile picture from form"})
		return
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log the error or handle it appropriately in a real application
			_ = closeErr
		}
	}()

	// Validate file type (only allow images)
	if !isValidImageType(handler.Header.Get("Content-Type")) {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid file type. Only images are allowed"})
		return
	}

	// Generate a unique filename
	filename := generateUniqueFilename(handler.Filename)

	// Create uploads directory if it doesn't exist
	err = os.MkdirAll("uploads/profile-pictures", os.ModePerm)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Unable to create uploads directory"})
		return
	}

	// Create the file
	dst, err := os.Create("uploads/profile-pictures/" + filename)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Unable to create file"})
		return
	}
	defer func() {
		if closeErr := dst.Close(); closeErr != nil {
			// Log the error or handle it appropriately in a real application
			_ = closeErr
		}
	}()

	// Copy the uploaded file to the destination
	_, err = io.Copy(dst, file)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Unable to save file"})
		return
	}

	// Update user's profile picture URL in the database
	profilePictureURL := "/uploads/profile-pictures/" + filename
	err = h.userService.UpdateProfilePictureURL(userID, profilePictureURL)
	if err != nil {
		// Try to delete the uploaded file since we couldn't update the database
		if err := os.Remove("uploads/profile-pictures/" + filename); err != nil {
			// Log the cleanup error, but don't send it to the client
			log.Printf("WARN: Failed to remove uploaded file %s during cleanup: %v", filename, err)
		}
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Unable to update profile picture"})
		return
	}

	// Return success response
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{
		"message":             "Profile picture uploaded successfully",
		"profile_picture_url": profilePictureURL,
	})
}

// isValidImageType checks if the content type is a valid image type
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}

	return false
}

// generateUniqueFilename generates a unique filename by adding a timestamp prefix
func generateUniqueFilename(originalFilename string) string {
	// Get file extension
	ext := filepath.Ext(originalFilename)

	// Generate timestamp
	timestamp := time.Now().UnixNano()

	// Generate random string
	randStr := generateRandomString(8)

	// Combine to create unique filename
	return fmt.Sprintf("%d_%s%s", timestamp, randStr, ext)
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		// Generate a random index
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Fallback to a simple approach if crypto/rand fails
			b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
			continue
		}
		b[i] = charset[n.Int64()]
	}

	return string(b)
}
