-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE friend_requests (
    id SERIAL PRIMARY KEY,
    from_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, accepted, rejected
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_friend_requests_from_user_id ON friend_requests(from_user_id);
CREATE INDEX idx_friend_requests_to_user_id ON friend_requests(to_user_id);
CREATE INDEX idx_friend_requests_status ON friend_requests(status);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE friend_requests;