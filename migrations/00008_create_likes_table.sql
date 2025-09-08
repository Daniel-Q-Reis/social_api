-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE likes (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    resource_type VARCHAR(50) NOT NULL, -- 'posts', 'photos', etc.
    resource_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, resource_type, resource_id)
);

CREATE INDEX idx_likes_resource ON likes(resource_type, resource_id);
CREATE INDEX idx_likes_user_id ON likes(user_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE likes;