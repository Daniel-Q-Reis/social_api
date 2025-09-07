-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    privacy VARCHAR(20) NOT NULL DEFAULT 'public', -- public, friends, only_me
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_created_at ON posts(created_at);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE posts;