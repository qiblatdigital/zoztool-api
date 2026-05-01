package main

import (
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/qiblatdigital/zoztool-api/internal/config"
	"github.com/qiblatdigital/zoztool-api/internal/handler"
	"github.com/qiblatdigital/zoztool-api/internal/middleware"
	"github.com/qiblatdigital/zoztool-api/internal/model"
	"github.com/qiblatdigital/zoztool-api/internal/repository"
	"github.com/qiblatdigital/zoztool-api/internal/service"
)

func main() {
	// Load .env
	_ = godotenv.Load()

	// Load config
	cfg := config.LoadConfig()
	cfg.Validate()

	// Connect DB
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(
		&model.User{},
		&model.SocialAccount{},
		&model.Post{},
		&model.PostTarget{},
		&model.MediaAsset{},
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	_ = repository.NewSocialAccountRepository(db)
	_ = repository.NewPostRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTRefreshSecret)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)

	// Setup router
	router := gin.New()
	router.Use(gin.Recovery())

	// Routes
	api := router.Group("/api/v1")
	{
		api.GET("/health", healthHandler.Health)

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh", authHandler.Refresh)
		}

		// Protected routes example (can add more later)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Add protected endpoints here
		}
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
