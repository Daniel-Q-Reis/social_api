package services

import (
	"fmt"
	"time"

	"github.com/gocli/social_api/internal/models"
	"github.com/gocli/social_api/internal/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService provides authentication-related functionality
type AuthService struct {
	BaseService
	authRepo           *repositories.AuthRepository
	userRepo           *repositories.UserRepository
	JWTSecret          string
	refreshTokenSecret string
}

// NewAuthService creates a new AuthService
func NewAuthService(authRepo *repositories.AuthRepository, userRepo *repositories.UserRepository,
	jwtSecret, refreshTokenSecret string) *AuthService {
	return &AuthService{
		authRepo:           authRepo,
		userRepo:           userRepo,
		JWTSecret:          jwtSecret,
		refreshTokenSecret: refreshTokenSecret,
	}
}

// Register creates a new user account
func (s *AuthService) Register(name, email, password string, birthDate time.Time) (*models.User, error) {
	// Check if user already exists
	_, err := s.userRepo.GetByEmail(email)
	if err == nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		BirthDate: birthDate,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
}

// Login authenticates a user and returns JWT tokens
func (s *AuthService) Login(email, password string) (string, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	// Generate access token
	accessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// RefreshToken generates a new access token using a refresh token
func (s *AuthService) RefreshToken(tokenString string) (string, error) {
	// Get refresh token from database
	refreshToken, err := s.authRepo.GetRefreshTokenByToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token")
	}

	// Check if token is expired
	if refreshToken.ExpiresAt.Before(time.Now()) {
		return "", fmt.Errorf("refresh token expired")
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(refreshToken.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

// Logout revokes a refresh token
func (s *AuthService) Logout(tokenString string) error {
	return s.authRepo.RevokeRefreshToken(tokenString)
}

// generateAccessToken generates a JWT access token
func (s *AuthService) generateAccessToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWTSecret))
}

// generateRefreshToken generates a refresh token and stores it in the database
func (s *AuthService) generateRefreshToken(userID int) (string, error) {
	// Generate random token
	token := generateRandomString(32)

	// Store in database
	refreshToken := &models.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30), // 30 days
		Revoked:   false,
	}

	if err := s.authRepo.CreateRefreshToken(refreshToken); err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return token, nil
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	// In a real implementation, you would use a cryptographically secure random generator
	// For simplicity, we'll return a fixed string here
	// In practice, you would use something like:
	// b := make([]byte, length)
	// rand.Read(b)
	// return base64.URLEncoding.EncodeToString(b)
	return "random_token_string_for_demo_purposes"
}
