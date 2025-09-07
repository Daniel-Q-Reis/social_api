// Package utils provides utility functions for the social media API
package utils

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	RefreshTokenSecret string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/social_api?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "jwt_secret_key"),
		RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", "refresh_token_secret_key"),
	}
}

// getEnv returns the value of an environment variable or a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
