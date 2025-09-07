package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gocli/social_api/internal/models"
)

// AuthRepository provides methods for accessing authentication data
type AuthRepository struct {
	*BaseRepository
}

// NewAuthRepository creates a new AuthRepository
func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{BaseRepository: NewBaseRepository(db)}
}

// CreateRefreshToken inserts a new refresh token into the database
func (r *AuthRepository) CreateRefreshToken(token *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at, revoked, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(query, token.UserID, token.Token, token.ExpiresAt,
		token.Revoked, now, now).Scan(&token.ID, &token.CreatedAt, &token.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}

	return nil
}

// GetRefreshTokenByToken retrieves a refresh token by its token value
func (r *AuthRepository) GetRefreshTokenByToken(token string) (*models.RefreshToken, error) {
	refreshToken := &models.RefreshToken{}
	query := `
		SELECT id, user_id, token, expires_at, revoked, created_at, updated_at
		FROM refresh_tokens
		WHERE token = $1 AND revoked = false`

	err := r.db.QueryRow(query, token).Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token,
		&refreshToken.ExpiresAt, &refreshToken.Revoked, &refreshToken.CreatedAt, &refreshToken.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found or revoked: %w", err)
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return refreshToken, nil
}

// RevokeRefreshToken revokes a refresh token
func (r *AuthRepository) RevokeRefreshToken(token string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true, updated_at = $1
		WHERE token = $2`

	now := time.Now()
	_, err := r.db.Exec(query, now, token)

	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	return nil
}
