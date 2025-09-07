-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE friends (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_friends_user_id ON friends(user_id);
CREATE INDEX idx_friends_friend_id ON friends(friend_id);
CREATE UNIQUE INDEX idx_friends_user_friend ON friends(user_id, friend_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE friends;