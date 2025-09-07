-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE albums (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    privacy VARCHAR(20) NOT NULL DEFAULT 'public', -- public, friends, only_me
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_albums_user_id ON albums(user_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE albums;