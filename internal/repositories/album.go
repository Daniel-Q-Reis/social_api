package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gocli/social_api/internal/models"
)

// AlbumRepository provides methods for accessing album data
type AlbumRepository struct {
	*BaseRepository
}

// NewAlbumRepository creates a new AlbumRepository
func NewAlbumRepository(db *sql.DB) *AlbumRepository {
	return &AlbumRepository{BaseRepository: NewBaseRepository(db)}
}

// CreateAlbum creates a new album
func (r *AlbumRepository) CreateAlbum(album *models.Album) error {
	query := `
		INSERT INTO albums (user_id, name, description, privacy, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(query, album.UserID, album.Name, album.Description, album.Privacy, now, now).
		Scan(&album.ID, &album.CreatedAt, &album.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create album: %w", err)
	}

	return nil
}

// GetAlbumByID retrieves an album by ID
func (r *AlbumRepository) GetAlbumByID(id int) (*models.Album, error) {
	album := &models.Album{}
	query := `
		SELECT id, user_id, name, description, privacy, created_at, updated_at
		FROM albums
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(&album.ID, &album.UserID, &album.Name,
		&album.Description, &album.Privacy, &album.CreatedAt, &album.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("album not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get album: %w", err)
	}

	return album, nil
}

// GetAlbumsByUserID retrieves albums for a specific user
func (r *AlbumRepository) GetAlbumsByUserID(userID int) ([]*models.Album, error) {
	albums := []*models.Album{}
	query := `
		SELECT id, user_id, name, description, privacy, created_at, updated_at
		FROM albums
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get albums: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// In a real application, you would log this error
			// For now, we'll just ignore it to keep the code simple
			_ = closeErr
		}
	}()

	for rows.Next() {
		album := &models.Album{}
		err := rows.Scan(&album.ID, &album.UserID, &album.Name, &album.Description,
			&album.Privacy, &album.CreatedAt, &album.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan album: %w", err)
		}
		albums = append(albums, album)
	}

	return albums, nil
}

// UpdateAlbum updates an album
func (r *AlbumRepository) UpdateAlbum(album *models.Album) error {
	query := `
		UPDATE albums
		SET name = $1, description = $2, privacy = $3, updated_at = $4
		WHERE id = $5`

	now := time.Now()
	_, err := r.db.Exec(query, album.Name, album.Description, album.Privacy, now, album.ID)

	if err != nil {
		return fmt.Errorf("failed to update album: %w", err)
	}

	album.UpdatedAt = now
	return nil
}

// DeleteAlbum deletes an album
func (r *AlbumRepository) DeleteAlbum(id int) error {
	// First delete all photos in the album
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			// In a real application, you would log this error
			// For now, we'll just ignore it to keep the code simple
			_ = rollbackErr
		}
	}()

	_, err = tx.Exec(`DELETE FROM photos WHERE album_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete photos: %w", err)
	}

	// Then delete the album
	_, err = tx.Exec(`DELETE FROM albums WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete album: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// CreatePhoto creates a new photo
func (r *AlbumRepository) CreatePhoto(photo *models.Photo) error {
	query := `
		INSERT INTO photos (album_id, url, caption, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(query, photo.AlbumID, photo.URL, photo.Caption, now, now).
		Scan(&photo.ID, &photo.CreatedAt, &photo.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create photo: %w", err)
	}

	return nil
}

// GetPhotosByAlbumID retrieves photos for a specific album
func (r *AlbumRepository) GetPhotosByAlbumID(albumID int) ([]*models.Photo, error) {
	photos := []*models.Photo{}
	query := `
		SELECT id, album_id, url, caption, created_at, updated_at
		FROM photos
		WHERE album_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, albumID)
	if err != nil {
		return nil, fmt.Errorf("failed to get photos: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			// In a real application, you would log this error
			// For now, we'll just ignore it to keep the code simple
			_ = closeErr
		}
	}()

	for rows.Next() {
		photo := &models.Photo{}
		err := rows.Scan(&photo.ID, &photo.AlbumID, &photo.URL, &photo.Caption, &photo.CreatedAt, &photo.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan photo: %w", err)
		}
		photos = append(photos, photo)
	}

	return photos, nil
}

// GetPhotoByID retrieves a photo by ID
func (r *AlbumRepository) GetPhotoByID(id int) (*models.Photo, error) {
	photo := &models.Photo{}
	query := `
		SELECT id, album_id, url, caption, created_at, updated_at
		FROM photos
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(&photo.ID, &photo.AlbumID, &photo.URL,
		&photo.Caption, &photo.CreatedAt, &photo.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("photo not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get photo: %w", err)
	}

	return photo, nil
}

// DeletePhoto deletes a photo
func (r *AlbumRepository) DeletePhoto(id int) error {
	query := `DELETE FROM photos WHERE id = $1`
	_, err := r.db.Exec(query, id)

	if err != nil {
		return fmt.Errorf("failed to delete photo: %w", err)
	}

	return nil
}
