# Social Media API

A complete, production-ready Social Media API built with Go following best practices.

## Features

- User authentication with JWT
- Profile management
- Friendships and connections
- Photo albums and media management
- News feed and posts
- Likes and comments

## Architecture

This project follows Clean Architecture principles with the following layers:

- `cmd/api` - Application entry point
- `internal/handlers` - HTTP handlers and routing
- `internal/services` - Business logic
- `internal/repositories` - Data access layer
- `internal/models` - Data structures and domain models
- `internal/middleware` - HTTP middleware
- `internal/database` - Database connection and setup
- `internal/utils` - Utility functions
- `migrations` - Database schema migrations

## Setup

1. Install dependencies:
   ```
   go mod tidy
   ```

2. Run the application:
   ```
   go run cmd/api/main.go
   ```

## Database Migrations

This project uses Goose for database migrations. To install Goose:

```
go install github.com/pressly/goose/v3/cmd/goose@latest
```

To run migrations:
```
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/social_api?sslmode=disable" up
```

## Testing

Run unit tests:
```
go test ./...
```

Run integration tests (requires Docker):
```
docker-compose up -d
go test -tags=integration ./...
```

## Linting

This project uses golangci-lint for code quality checks.

To install golangci-lint:
```
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

Run the linter:
```
golangci-lint run
```

## Docker

To build and run the application with Docker:
```
docker-compose up --build
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and get JWT tokens
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout and revoke refresh token

### Users
- `GET /api/v1/users/{userId}` - Get a user's public profile
- `GET /api/v1/users/search?q={query}` - Search for users
- `GET /api/v1/me` - Get the logged-in user's full profile
- `PUT /api/v1/me` - Full update of the logged-in user's profile
- `PATCH /api/v1/me` - Partial update of the logged-in user's profile

### Friendships
- `GET /api/v1/users/{userId}/friends` - List a user's friends
- `GET /api/v1/me/friend-requests` - List pending friend requests for the logged-in user
- `POST /api/v1/users/{userId}/friend-requests` - Send a friend request
- `POST /api/v1/friend-requests/{requestId}/accept` - Accept a friend request
- `POST /api/v1/friend-requests/{requestId}/reject` - Reject a friend request
- `DELETE /api/v1/users/{userId}/friends` - Unfriend a user

### Albums and Photos
- `POST /api/v1/me/albums` - Create a new photo album
- `GET /api/v1/users/{userId}/albums` - List a user's albums
- `GET /api/v1/albums/{albumId}` - Get album details
- `PUT /api/v1/albums/{albumId}` - Update album info
- `DELETE /api/v1/albums/{albumId}` - Delete an album

### Feed and Posts
- `POST /api/v1/posts` - Create a new post
- `GET /api/v1/feed` - Get the personalized news feed
- `GET /api/v1/users/{userId}/posts` - List a user's posts
- `GET /api/v1/posts/{postId}` - Get a single post
- `PUT /api/v1/posts/{postId}` - Edit a post
- `DELETE /api/v1/posts/{postId}` - Delete a post