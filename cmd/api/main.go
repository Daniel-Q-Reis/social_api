package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/gocli/social_api/internal/database"
	"github.com/gocli/social_api/internal/handlers"
	middlewares "github.com/gocli/social_api/internal/middleware"
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
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.CORSMiddleware) // Add CORS middleware

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.SendJSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Auth routes
	router.Post("/api/v1/auth/register", authHandler.Register)
	router.Post("/api/v1/auth/login", authHandler.Login)
	router.Post("/api/v1/auth/refresh", authHandler.RefreshToken)
	router.Post("/api/v1/auth/logout", authHandler.Logout)

	// Protected routes
	router.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(authService))

		// User routes
		r.Get("/api/v1/users/{userId}", userHandler.GetUserProfile)
		r.Get("/api/v1/users/search", userHandler.SearchUsers)
		r.Get("/api/v1/me", userHandler.GetMe)
		r.Put("/api/v1/me", userHandler.UpdateMe)
		r.Patch("/api/v1/me", userHandler.PartialUpdateMe)
		r.Post("/api/v1/me/profile-picture", userHandler.UploadProfilePicture)

		// Friend routes
		r.Get("/api/v1/users/{userId}/friends", friendHandler.GetUserFriends)
		r.Get("/api/v1/me/friend-requests", friendHandler.GetMyFriendRequests)
		r.Post("/api/v1/users/{userId}/friend-requests", friendHandler.SendFriendRequest)
		r.Post("/api/v1/friend-requests/{requestId}/accept", friendHandler.AcceptFriendRequest)
		r.Post("/api/v1/friend-requests/{requestId}/reject", friendHandler.RejectFriendRequest)
		r.Delete("/api/v1/users/{userId}/friends", friendHandler.UnfriendUser)

		// Post routes
		r.Post("/api/v1/posts", postHandler.CreatePost)
		r.Get("/api/v1/feed", postHandler.GetFeed)
		r.Get("/api/v1/users/{userId}/posts", postHandler.GetUserPosts)
		r.Get("/api/v1/posts/{postId}", postHandler.GetPost)
		r.Put("/api/v1/posts/{postId}", postHandler.UpdatePost)
		r.Delete("/api/v1/posts/{postId}", postHandler.DeletePost)

		// Album routes
		r.Post("/api/v1/me/albums", albumHandler.CreateAlbum)
		r.Get("/api/v1/users/{userId}/albums", albumHandler.GetUserAlbums)
		r.Get("/api/v1/albums/{albumId}", albumHandler.GetAlbum)
		r.Put("/api/v1/albums/{albumId}", albumHandler.UpdateAlbum)
		r.Delete("/api/v1/albums/{albumId}", albumHandler.DeleteAlbum)

		// Like routes
		r.Post("/api/v1/{resourceType}/{resourceId}/like", likeHandler.LikeResource)
		r.Delete("/api/v1/{resourceType}/{resourceId}/like", likeHandler.UnlikeResource)
		r.Get("/api/v1/{resourceType}/{resourceId}/likes", likeHandler.GetLikesForResource)

		// Comment routes
		r.Post("/api/v1/{resourceType}/{resourceId}/comments", commentHandler.CreateComment)
		r.Get("/api/v1/{resourceType}/{resourceId}/comments", commentHandler.GetCommentsForResource)
		r.Delete("/api/v1/comments/{commentId}", commentHandler.DeleteComment)
	})

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
