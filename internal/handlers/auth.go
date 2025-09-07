package handlers

import (
	"net/http"
	"time"

	"github.com/gocli/social_api/internal/services"
	"github.com/gocli/social_api/internal/utils"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	BaseHandler
	authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req struct {
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Password  string    `json:"password"`
		BirthDate time.Time `json:"birth_date"`
	}

	if err := h.DecodeJSONBody(w, r, &req); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	// Register user
	user, err := h.authService.Register(req.Name, req.Email, req.Password, req.BirthDate)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Return user data
	utils.SendJSONResponse(w, http.StatusCreated, user)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := h.DecodeJSONBody(w, r, &req); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	// Authenticate user
	accessToken, refreshToken, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	// Return tokens
	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// RefreshToken handles refresh token requests
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := h.DecodeJSONBody(w, r, &req); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return
	}

	// Refresh token
	accessToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	// Return new access token
	response := map[string]string{
		"access_token": accessToken,
	}
	utils.SendJSONResponse(w, http.StatusOK, response)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get the refresh token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		utils.SendJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Refresh token required"})
		return
	}

	refreshToken := authHeader[7:] // Remove "Bearer " prefix

	// Logout user
	if err := h.authService.Logout(refreshToken); err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Return success response
	utils.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}
