// Package handlers provides HTTP handlers for the social media API
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gocli/social_api/internal/middleware"
)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	// Add common handler functionality here if needed
}

// GetUserIDFromContext extracts the user ID from the request context
func (h *BaseHandler) GetUserIDFromContext(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(middleware.UserContextKey).(int)
	if !ok {
		return 0, http.ErrNoCookie
	}
	return userID, nil
}

// ParseIDFromPath parses an ID from the URL path
func (h *BaseHandler) ParseIDFromPath(r *http.Request, param string) (int, error) {
	// In a real implementation, you would use a router that supports path parameters
	// For simplicity, we'll return a fixed value here
	// In practice, you would use something like:
	// vars := mux.Vars(r)
	// idStr := vars[param]
	idStr := "1" // Placeholder
	return strconv.Atoi(idStr)
}

// ParseQueryInt parses an integer from query parameters
func (h *BaseHandler) ParseQueryInt(r *http.Request, param string, defaultValue int) int {
	valueStr := r.URL.Query().Get(param)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// DecodeJSONBody decodes JSON from the request body
func (h *BaseHandler) DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}
