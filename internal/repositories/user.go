package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gocli/social_api/internal/models"
)

// UserRepository provides methods for accessing user data
type UserRepository struct {
	*BaseRepository
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{BaseRepository: NewBaseRepository(db)}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (name, email, password, birth_date, profile_picture_url, cover_photo_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(query, user.Name, user.Email, user.Password, user.BirthDate,
		user.ProfilePictureURL, user.CoverPhotoURL, now, now).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.Password = "" // Clear password from memory
	return nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, password, birth_date, profile_picture_url, cover_photo_url, created_at, updated_at
		FROM users
		WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password,
		&user.BirthDate, &user.ProfilePictureURL, &user.CoverPhotoURL, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, birth_date, profile_picture_url, cover_photo_url, created_at, updated_at
		FROM users
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email,
		&user.BirthDate, &user.ProfilePictureURL, &user.CoverPhotoURL, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

// Update updates a user's information
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, birth_date = $3, profile_picture_url = $4, cover_photo_url = $5, updated_at = $6
		WHERE id = $7`

	now := time.Now()
	_, err := r.db.Exec(query, user.Name, user.Email, user.BirthDate,
		user.ProfilePictureURL, user.CoverPhotoURL, now, user.ID)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	user.UpdatedAt = now
	return nil
}

// Search searches for users by name or email
func (r *UserRepository) Search(query string, limit, offset int) ([]*models.UserPublic, error) {
	users := []*models.UserPublic{}
	sqlQuery := `
		SELECT id, name, profile_picture_url, cover_photo_url, created_at
		FROM users
		WHERE name ILIKE $1 OR email ILIKE $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(sqlQuery, "%"+query+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// In a real application, you would log this error
			// For now, we'll just ignore it to keep the code simple
			_ = closeErr
		}
	}()

	for rows.Next() {
		user := &models.UserPublic{}
		err := rows.Scan(&user.ID, &user.Name, &user.ProfilePictureURL, &user.CoverPhotoURL, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}
