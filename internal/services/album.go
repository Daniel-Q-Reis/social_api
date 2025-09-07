package services

import (
	"fmt"

	"github.com/gocli/social_api/internal/models"
	"github.com/gocli/social_api/internal/repositories"
)

// AlbumService provides album-related functionality
type AlbumService struct {
	BaseService
	albumRepo *repositories.AlbumRepository
}

// NewAlbumService creates a new AlbumService
func NewAlbumService(albumRepo *repositories.AlbumRepository) *AlbumService {
	return &AlbumService{
		albumRepo: albumRepo,
	}
}

// CreateAlbum creates a new album
func (s *AlbumService) CreateAlbum(userID int, name, description, privacy string) (*models.Album, error) {
	album := &models.Album{
		UserID:      userID,
		Name:        name,
		Description: description,
		Privacy:     privacy,
	}

	if err := s.albumRepo.CreateAlbum(album); err != nil {
		return nil, fmt.Errorf("failed to create album: %w", err)
	}

	return album, nil
}

// GetAlbumByID retrieves an album by ID
func (s *AlbumService) GetAlbumByID(id int) (*models.Album, error) {
	album, err := s.albumRepo.GetAlbumByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get album: %w", err)
	}
	return album, nil
}

// GetAlbumsWithPhotosByUserID retrieves albums with their photos for a specific user
func (s *AlbumService) GetAlbumsWithPhotosByUserID(userID int) ([]*models.Album, [][]*models.Photo, error) {
	albums, err := s.albumRepo.GetAlbumsByUserID(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get albums: %w", err)
	}

	// Get photos for each album
	albumsPhotos := make([][]*models.Photo, len(albums))
	for i, album := range albums {
		photos, err := s.albumRepo.GetPhotosByAlbumID(album.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get photos for album %d: %w", album.ID, err)
		}
		albumsPhotos[i] = photos
	}

	return albums, albumsPhotos, nil
}

// UpdateAlbum updates an album
func (s *AlbumService) UpdateAlbum(albumID, userID int, name, description, privacy string) (*models.Album, error) {
	// First get the album to verify ownership
	album, err := s.albumRepo.GetAlbumByID(albumID)
	if err != nil {
		return nil, fmt.Errorf("failed to get album: %w", err)
	}

	// Check if user is authorized to update this album
	if album.UserID != userID {
		return nil, fmt.Errorf("user is not authorized to update this album")
	}

	// Update the album
	album.Name = name
	album.Description = description
	album.Privacy = privacy

	if err := s.albumRepo.UpdateAlbum(album); err != nil {
		return nil, fmt.Errorf("failed to update album: %w", err)
	}

	return album, nil
}

// DeleteAlbum deletes an album
func (s *AlbumService) DeleteAlbum(albumID, userID int) error {
	// First get the album to verify ownership
	album, err := s.albumRepo.GetAlbumByID(albumID)
	if err != nil {
		return fmt.Errorf("failed to get album: %w", err)
	}

	// Check if user is authorized to delete this album
	if album.UserID != userID {
		return fmt.Errorf("user is not authorized to delete this album")
	}

	// Delete the album
	if err := s.albumRepo.DeleteAlbum(albumID); err != nil {
		return fmt.Errorf("failed to delete album: %w", err)
	}

	return nil
}

// AddPhotoToAlbum adds a photo to an album
func (s *AlbumService) AddPhotoToAlbum(albumID, userID int, url, caption string) (*models.Photo, error) {
	// First get the album to verify ownership
	album, err := s.albumRepo.GetAlbumByID(albumID)
	if err != nil {
		return nil, fmt.Errorf("failed to get album: %w", err)
	}

	// Check if user is authorized to add photos to this album
	if album.UserID != userID {
		return nil, fmt.Errorf("user is not authorized to add photos to this album")
	}

	// Create the photo
	photo := &models.Photo{
		AlbumID: albumID,
		URL:     url,
		Caption: caption,
	}

	if err := s.albumRepo.CreatePhoto(photo); err != nil {
		return nil, fmt.Errorf("failed to create photo: %w", err)
	}

	return photo, nil
}

// DeletePhoto deletes a photo
func (s *AlbumService) DeletePhoto(photoID, userID int) error {
	// First get the photo to verify ownership
	photo, err := s.albumRepo.GetPhotoByID(photoID)
	if err != nil {
		return fmt.Errorf("failed to get photo: %w", err)
	}

	// Get the album to verify user ownership
	album, err := s.albumRepo.GetAlbumByID(photo.AlbumID)
	if err != nil {
		return fmt.Errorf("failed to get album: %w", err)
	}

	// Check if user is authorized to delete this photo
	if album.UserID != userID {
		return fmt.Errorf("user is not authorized to delete this photo")
	}

	// Delete the photo
	if err := s.albumRepo.DeletePhoto(photoID); err != nil {
		return fmt.Errorf("failed to delete photo: %w", err)
	}

	return nil
}
