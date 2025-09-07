-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    birth_date DATE NOT NULL,
    profile_picture_url VARCHAR(500),
    cover_photo_url VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE users;