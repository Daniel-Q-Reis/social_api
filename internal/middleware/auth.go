// Package middleware provides HTTP middleware for the social media API
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gocli/social_api/internal/services"
	"github.com/gocli/social_api/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

// ContextKey is a custom type for context keys
type ContextKey string

// UserContextKey is the key for storing user ID in context
const UserContextKey ContextKey = "userID"

// AuthMiddleware is a middleware that verifies JWT tokens
func AuthMiddleware(next http.Handler, authService *services.AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Authorization header required"})
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid authorization header format"})
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(authService.JWTSecret), nil
		})

		if err != nil {
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			return
		}

		// Check if the token is valid
		if !token.Valid {
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			return
		}

		// Extract user ID from claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			utils.SendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid user ID in token"})
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserContextKey, int(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
