# Social Media API

Welcome to the Social Media API, a robust, production-ready backend service built with Go. This project showcases a complete set of features for a modern social media platform, architected with best practices and clean code principles.

## 🚀 Features

* **User Authentication**: Secure user registration and login using JWT (Access and Refresh Tokens).
* **Profile Management**: Full control over user profiles, including profile picture uploads.
* **Social Graph**: Functionality for friendships, including sending, accepting, and rejecting friend requests.
* **Content Creation**: Users can create, edit, and delete posts with different privacy levels (public, friends, only me).
* **News Feed**: A personalized news feed to view posts from friends.
* **Media Management**: Support for photo albums and media uploads.
* **Engagement**: Interactive features like likes and comments on posts and other resources.

## 🏛️ Architecture

This project is built upon the principles of **Clean Architecture** to ensure a scalable, maintainable, and testable codebase. The project is organized into the following layers:

* `cmd/api`: The application's entry point, responsible for bootstrapping the server.
* `internal/handlers`: HTTP handlers that manage requests and responses.
* `internal/services`: The core business logic of the application resides here.
* `internal/repositories`: The data access layer, responsible for database interactions.
* `internal/models`: Defines the data structures and domain models.
* `internal/middleware`: Custom HTTP middleware, including authentication.
* `internal/database`: Manages the database connection and setup.
* `internal/utils`: Contains utility functions for configuration, validation, and responses.
* `migrations`: Database schema migrations managed by Goose.

## 🛠️ Getting Started

Follow these instructions to get the project up and running on your local machine.

### Prerequisites

* Go (version 1.24 or later)
* Docker and Docker Compose
* [Goose](https://github.com/pressly/goose) for database migrations

### Installation & Setup

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/daniel-q-reis/social_api.git](https://github.com/daniel-q-reis/social_api.git)
    cd social_api
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Run with Docker (Recommended):**
    The easiest way to get started is by using Docker Compose, which will set up the database and run the application.
    ```bash
    docker-compose up --build
    ```
    The API will be available at `http://localhost:8080`.

4.  **Running Locally (Without Docker):**
    If you prefer to run the application and database manually:

    * **Start the PostgreSQL database:** You can use Docker to easily spin up a Postgres instance.
        ```bash
        docker run --name social-db -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=social_api -p 5432:5432 -d postgres:13
        ```

    * **Run database migrations:** This project uses Goose for managing database schema.
        ```bash
        # Install Goose if you haven't already
        go install [github.com/pressly/goose/v3/cmd/goose@latest](https://github.com/pressly/goose/v3/cmd/goose@latest)

        # Run migrations
        goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/social_api?sslmode=disable" up
        ```

    * **Run the application:**
        ```bash
        go run cmd/api/main.go
        ```

## 🧪 Testing and Linting

### Running Tests

* **Unit Tests:**
    ```bash
    go test ./...
    ```

* **Integration Tests:** These tests require a running database. Ensure the test database is up (e.g., via `docker-compose up -d db`).
    ```bash
    go test -tags=integration ./...
    ```

### Linting

This project uses `golangci-lint` for ensuring code quality.

* **Install linter:**
    ```bash
    go install [github.com/golangci/golangci-lint/cmd/golangci-lint@latest](https://github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
    ```

* **Run linter:**
    ```bash
    golangci-lint run
    ```

## 📜 API Endpoints

All endpoints are prefixed with `/api/v1`.

### Authentication

| Method | Endpoint              | Description                      |
| :----- | :-------------------- | :------------------------------- |
| `POST` | `/auth/register`      | Register a new user              |
| `POST` | `/auth/login`         | Login and get JWT tokens         |
| `POST` | `/auth/refresh`       | Refresh access token             |
| `POST` | `/auth/logout`        | Logout and revoke refresh token  |

### Users

| Method  | Endpoint                  | Description                                |
| :------ | :------------------------ | :----------------------------------------- |
| `GET`   | `/users/{userId}`         | Get a user's public profile                |
| `GET`   | `/users/search?q={query}` | Search for users by name or email          |
| `GET`   | `/me`                     | Get the logged-in user's full profile      |
| `PUT`   | `/me`                     | Full update of the logged-in user's profile |
| `PATCH` | `/me`                     | Partial update of the user's profile       |
| `POST`  | `/me/profile-picture`     | Upload a profile picture for the user      |

### Friendships

| Method   | Endpoint                               | Description                                      |
| :------- | :------------------------------------- | :----------------------------------------------- |
| `GET`    | `/users/{userId}/friends`              | List a user's friends                            |
| `GET`    | `/me/friend-requests`                  | List pending friend requests for the logged-in user |
| `POST`   | `/users/{userId}/friend-requests`      | Send a friend request to a user                  |
| `POST`   | `/friend-requests/{requestId}/accept`  | Accept a pending friend request                  |
| `POST`   | `/friend-requests/{requestId}/reject`  | Reject a pending friend request                  |
| `DELETE` | `/users/{userId}/friends`              | Unfriend a user                                  |

### Posts & Feed

| Method   | Endpoint                  | Description                     |
| :------- | :------------------------ | :------------------------------ |
| `POST`   | `/posts`                  | Create a new post               |
| `GET`    | `/feed`                   | Get the personalized news feed  |
| `GET`    | `/users/{userId}/posts`   | List a user's posts             |
| `GET`    | `/posts/{postId}`         | Get a single post               |
| `PUT`    | `/posts/{postId}`         | Edit an existing post           |
| `DELETE` | `/posts/{postId}`         | Delete a post                   |

### Albums & Photos

| Method   | Endpoint                   | Description                      |
| :------- | :------------------------- | :------------------------------- |
| `POST`   | `/me/albums`               | Create a new photo album         |
| `GET`    | `/users/{userId}/albums`   | List a user's photo albums       |
| `GET`    | `/albums/{albumId}`        | Get details for a single album   |
| `PUT`    | `/albums/{albumId}`        | Update album information         |
| `DELETE` | `/albums/{albumId}`        | Delete an album and its photos   |

### Likes & Comments

| Method   | Endpoint                                 | Description                        |
| :------- | :--------------------------------------- | :--------------------------------- |
| `POST`   | `/{resourceType}/{resourceId}/like`      | Like a resource (e.g., post, photo) |
| `DELETE` | `/{resourceType}/{resourceId}/like`      | Unlike a resource                  |
| `GET`    | `/{resourceType}/{resourceId}/likes`     | Get likes for a resource           |
| `POST`   | `/{resourceType}/{resourceId}/comments`  | Add a comment to a resource        |
| `GET`    | `/{resourceType}/{resourceId}/comments`  | Get comments for a resource        |
| `DELETE` | `/comments/{commentId}`                  | Delete a comment                   |

## 🙏 Acknowledgments

This project was developed with the significant use of AI-powered tools and technologies. Generative AI was instrumental in accelerating development, generating boilerplate code, writing tests, and providing architectural insights, leading to a more efficient and robust development process.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.