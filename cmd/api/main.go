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
	"github.com/qiblatdigital/zoztool-api/internal/pkg/threads"
	"github.com/qiblatdigital/zoztool-api/internal/repository"
	"github.com/qiblatdigital/zoztool-api/internal/scheduler"
	"github.com/qiblatdigital/zoztool-api/internal/service"
)

func main() {
	_ = godotenv.Load()

	cfg := config.LoadConfig()
	cfg.Validate()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	socialAccountRepo := repository.NewSocialAccountRepository(db)
	postRepo := repository.NewPostRepository(db)
	postTargetRepo := repository.NewPostTargetRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTRefreshSecret)
	socialAccountService := service.NewSocialAccountService(socialAccountRepo)
	threadsClient := threads.NewClient()
	postService := service.NewPostService(postRepo, postTargetRepo, socialAccountRepo, threadsClient)

	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)
	socialAccountHandler := handler.NewSocialAccountHandler(socialAccountService)
	postHandler := handler.NewPostHandler(postService)

	sched := scheduler.NewScheduler(postService)
	sched.Start()

	router := gin.New()
	router.Use(gin.Recovery())

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

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			sa := protected.Group("/social-accounts")
			{
				sa.POST("", socialAccountHandler.Create)
				sa.GET("", socialAccountHandler.List)
				sa.GET("/:id", socialAccountHandler.GetByID)
				sa.DELETE("/:id", socialAccountHandler.Delete)
			}

			posts := protected.Group("/posts")
			{
				posts.POST("", postHandler.Create)
				posts.GET("", postHandler.List)
				posts.GET("/:id", postHandler.GetByID)
				posts.PUT("/:id", postHandler.Update)
				posts.DELETE("/:id", postHandler.Delete)
				posts.POST("/:id/publish", postHandler.Publish)
			}
		}
	}

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
