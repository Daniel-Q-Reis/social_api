package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gocli/social_api/internal/database"
	"github.com/gocli/social_api/internal/handlers"
	"github.com/gocli/social_api/internal/middleware"
	"github.com/gocli/social_api/internal/repositories"
	"github.com/gocli/social_api/internal/services"
	"github.com/gocli/social_api/internal/utils"
)

func main() {
	// Initialize configuration
	config := utils.LoadConfig()

	// Initialize validator
	validator := utils.NewValidator()

	// Initialize database connection
	db, err := database.Connect(config.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database connection: %v", closeErr)
		}
	}()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	authRepo := repositories.NewAuthRepository(db)
	friendRepo := repositories.NewFriendRepository(db)
	postRepo := repositories.NewPostRepository(db)
	albumRepo := repositories.NewAlbumRepository(db)
	likeRepo := repositories.NewLikeRepository(db)
	commentRepo := repositories.NewCommentRepository(db)

	// Initialize services
	authService := services.NewAuthService(authRepo, userRepo, config.JWTSecret, config.RefreshTokenSecret)
	userService := services.NewUserService(userRepo)
	friendService := services.NewFriendService(friendRepo, userRepo)
	postService := services.NewPostService(postRepo)
	albumService := services.NewAlbumService(albumRepo)
	likeService := services.NewLikeService(likeRepo)
	commentService := services.NewCommentService(commentRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, validator)
	userHandler := handlers.NewUserHandler(userService, validator)
	friendHandler := handlers.NewFriendHandler(friendService, validator)
	postHandler := handlers.NewPostHandler(postService, validator)
	albumHandler := handlers.NewAlbumHandler(albumService, validator)
	likeHandler := handlers.NewLikeHandler(likeService, validator)
	commentHandler := handlers.NewCommentHandler(commentService, validator)

	// Setup routes
	router := http.NewServeMux()

	// Health check endpoint
	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Auth routes
	router.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	router.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	router.HandleFunc("POST /api/v1/auth/refresh", authHandler.RefreshToken)
	router.HandleFunc("POST /api/v1/auth/logout", authHandler.Logout)

	// Protected routes
	protectedRouter := http.NewServeMux()

	// User routes
	protectedRouter.HandleFunc("GET /api/v1/users/{userId}", userHandler.GetUserProfile)
	protectedRouter.HandleFunc("GET /api/v1/users/search", userHandler.SearchUsers)
	protectedRouter.HandleFunc("GET /api/v1/me", userHandler.GetMe)
	protectedRouter.HandleFunc("PUT /api/v1/me", userHandler.UpdateMe)
	protectedRouter.HandleFunc("PATCH /api/v1/me", userHandler.PartialUpdateMe)

	// Friend routes
	protectedRouter.HandleFunc("GET /api/v1/users/{userId}/friends", friendHandler.GetUserFriends)
	protectedRouter.HandleFunc("GET /api/v1/me/friend-requests", friendHandler.GetMyFriendRequests)
	protectedRouter.HandleFunc("POST /api/v1/users/{userId}/friend-requests", friendHandler.SendFriendRequest)
	protectedRouter.HandleFunc("POST /api/v1/friend-requests/{requestId}/accept", friendHandler.AcceptFriendRequest)
	protectedRouter.HandleFunc("POST /api/v1/friend-requests/{requestId}/reject", friendHandler.RejectFriendRequest)
	protectedRouter.HandleFunc("DELETE /api/v1/users/{userId}/friends", friendHandler.UnfriendUser)

	// Post routes
	protectedRouter.HandleFunc("POST /api/v1/posts", postHandler.CreatePost)
	protectedRouter.HandleFunc("GET /api/v1/feed", postHandler.GetFeed)
	protectedRouter.HandleFunc("GET /api/v1/users/{userId}/posts", postHandler.GetUserPosts)
	protectedRouter.HandleFunc("GET /api/v1/posts/{postId}", postHandler.GetPost)
	protectedRouter.HandleFunc("PUT /api/v1/posts/{postId}", postHandler.UpdatePost)
	protectedRouter.HandleFunc("DELETE /api/v1/posts/{postId}", postHandler.DeletePost)

	// Album routes
	protectedRouter.HandleFunc("POST /api/v1/me/albums", albumHandler.CreateAlbum)
	protectedRouter.HandleFunc("GET /api/v1/users/{userId}/albums", albumHandler.GetUserAlbums)
	protectedRouter.HandleFunc("GET /api/v1/albums/{albumId}", albumHandler.GetAlbum)
	protectedRouter.HandleFunc("PUT /api/v1/albums/{albumId}", albumHandler.UpdateAlbum)
	protectedRouter.HandleFunc("DELETE /api/v1/albums/{albumId}", albumHandler.DeleteAlbum)

	// Like routes
	protectedRouter.HandleFunc("POST /api/v1/{resourceType}/{resourceId}/like", likeHandler.LikeResource)
	protectedRouter.HandleFunc("DELETE /api/v1/{resourceType}/{resourceId}/like", likeHandler.UnlikeResource)
	protectedRouter.HandleFunc("GET /api/v1/{resourceType}/{resourceId}/likes", likeHandler.GetLikesForResource)

	// Comment routes
	protectedRouter.HandleFunc("POST /api/v1/{resourceType}/{resourceId}/comments", commentHandler.CreateComment)
	protectedRouter.HandleFunc("GET /api/v1/{resourceType}/{resourceId}/comments", commentHandler.GetCommentsForResource)
	protectedRouter.HandleFunc("DELETE /api/v1/comments/{commentId}", commentHandler.DeleteComment)

	// Apply middleware to protected routes
	protectedHandler := middleware.AuthMiddleware(protectedRouter, authService)

	// Mount protected routes
	router.Handle("/", protectedHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
