-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE photos (
    id SERIAL PRIMARY KEY,
    album_id INTEGER NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
    url VARCHAR(500) NOT NULL,
    caption TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_photos_album_id ON photos(album_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE photos;