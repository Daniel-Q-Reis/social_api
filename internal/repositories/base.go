// Package repositories provides data access functionality for the social media API
package repositories

import (
	"database/sql"
)

// BaseRepository provides common functionality for all repositories
type BaseRepository struct {
	db *sql.DB
}

// NewBaseRepository creates a new BaseRepository
func NewBaseRepository(db *sql.DB) *BaseRepository {
	return &BaseRepository{db: db}
}
