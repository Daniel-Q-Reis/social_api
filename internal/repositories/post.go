package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gocli/social_api/internal/models"
)

// PostRepository provides methods for accessing post data
type PostRepository struct {
	*BaseRepository
}

// NewPostRepository creates a new PostRepository
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{BaseRepository: NewBaseRepository(db)}
}

// Create inserts a new post into the database
func (r *PostRepository) Create(post *models.Post) error {
	query := `
		INSERT INTO posts (user_id, content, privacy, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(query, post.UserID, post.Content, post.Privacy, now, now).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

// GetByID retrieves a post by ID
func (r *PostRepository) GetByID(id int) (*models.Post, error) {
	post := &models.Post{}
	query := `
		SELECT id, user_id, content, privacy, created_at, updated_at
		FROM posts
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(&post.ID, &post.UserID, &post.Content,
		&post.Privacy, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return post, nil
}

// GetPostsByUserID retrieves posts for a specific user
func (r *PostRepository) GetPostsByUserID(userID int, limit, offset int) ([]*models.Post, error) {
	posts := []*models.Post{}
	query := `
		SELECT id, user_id, content, privacy, created_at, updated_at
		FROM posts
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// In a real application, you would log this error
			// For now, we'll just ignore it to keep the code simple
			_ = closeErr
		}
	}()

	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Privacy, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// GetFeed retrieves posts for a user's feed
func (r *PostRepository) GetFeed(userID int, limit, offset int) ([]*models.PostWithUser, error) {
	posts := []*models.PostWithUser{}
	query := `
		SELECT p.id, p.user_id, p.content, p.privacy, p.created_at, p.updated_at,
		       u.name, u.profile_picture_url
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id IN (
		    SELECT friend_id FROM friends WHERE user_id = $1
		    UNION
		    SELECT $1
		)
		AND (p.privacy = 'public' OR p.privacy = 'friends' OR (p.privacy = 'only_me' AND p.user_id = $1))
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// In a real application, you would log this error
			// For now, we'll just ignore it to keep the code simple
			_ = closeErr
		}
	}()

	for rows.Next() {
		post := &models.PostWithUser{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Privacy,
			&post.CreatedAt, &post.UpdatedAt, &post.UserName, &post.ProfilePictureURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// Update updates a post
func (r *PostRepository) Update(post *models.Post) error {
	query := `
		UPDATE posts
		SET content = $1, privacy = $2, updated_at = $3
		WHERE id = $4`

	now := time.Now()
	_, err := r.db.Exec(query, post.Content, post.Privacy, now, post.ID)

	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	post.UpdatedAt = now
	return nil
}

// Delete deletes a post
func (r *PostRepository) Delete(id int) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.Exec(query, id)

	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}
