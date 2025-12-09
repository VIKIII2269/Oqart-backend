package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oqart/backend/config"
	"github.com/oqart/backend/internal/delivery/http/handler"
	"github.com/oqart/backend/internal/delivery/http/middleware"
	"github.com/oqart/backend/internal/repository/postgres"
	"github.com/oqart/backend/internal/usecase"
	"github.com/oqart/backend/pkg/cache"
	"github.com/oqart/backend/pkg/database"
	"github.com/oqart/backend/pkg/jwt"
	"github.com/oqart/backend/pkg/logger"
)

// @title OQart Backend API
// @version 1.0
// @description Production-grade e-commerce backend for OQart - India's premier organic products marketplace
// @termsOfService https://oqart.com/terms

// @contact.name API Support
// @contact.url https://oqart.com/support
// @contact.email support@oqart.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.InitGlobal(cfg.Logging.Level, cfg.Logging.Format, cfg.Logging.Output)
	log := logger.Global

	log.Info("Starting OQart Backend API",
		"version", cfg.Server.Version,
		"env", cfg.Server.Env)

	// Initialize database
	db, err := database.NewPostgres(&cfg.Database, log)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	defer db.Close()

	log.Info("Database connection established")

	// Initialize Redis
	var redisCache *cache.Cache
	redisCache, err = cache.NewRedis(&cfg.Redis, log)
	if err != nil {
		log.Warn("Failed to connect to Redis, continuing without cache", err)
		redisCache = nil
	} else {
		defer redisCache.Close()
		log.Info("Redis connection established")
	}

	// Initialize JWT manager
	jwtManager := jwt.NewTokenManager(&cfg.JWT)

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db.DB)
	sessionRepo := postgres.NewSessionRepository(db.DB)
	otpRepo := postgres.NewOTPRepository(db.DB)
	passwordResetRepo := postgres.NewPasswordResetRepository(db.DB)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(
		userRepo,
		sessionRepo,
		otpRepo,
		passwordResetRepo,
		jwtManager,
		redisCache,
		cfg,
		log,
	)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUseCase, log)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtManager, userRepo)

	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.Default()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(corsMiddleware(cfg))

	// Health check endpoint
	router.GET("/health", healthCheckHandler(db, redisCache))
	router.GET("/api/v1/health", healthCheckHandler(db, redisCache))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
				"version": cfg.Server.Version,
				"env":     cfg.Server.Env,
			})
		})

		// Authentication routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/login/phone", authHandler.SendPhoneOTP)
			auth.POST("/verify-otp", authHandler.VerifyPhoneOTP)
			auth.POST("/refresh-token", authHandler.RefreshToken)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.GET("/verify-email", authHandler.VerifyEmail)

			// Protected auth routes
			authProtected := auth.Group("")
			authProtected.Use(authMiddleware.RequireAuth())
			{
				authProtected.POST("/logout", authHandler.Logout)
			}
		}

		// User routes (protected)
		users := v1.Group("/users")
		users.Use(authMiddleware.RequireAuth())
		{
			// TODO: Add user management endpoints
			// users.GET("/me", userHandler.GetMe)
			// users.PUT("/me", userHandler.UpdateMe)
			// users.PUT("/me/password", userHandler.ChangePassword)
			// users.PUT("/me/email", userHandler.UpdateEmail)
			// users.PUT("/me/phone", userHandler.UpdatePhone)
			// users.DELETE("/me", userHandler.DeleteAccount)
			// users.GET("/me/preferences", userHandler.GetPreferences)
			// users.PUT("/me/preferences", userHandler.UpdatePreferences)
		}

		// Vendor routes
		vendors := v1.Group("/vendors")
		{
			// TODO: Add vendor endpoints
			// vendors.POST("/onboard/step1", vendorHandler.OnboardStep1)
			// vendors.POST("/verify-gstin", vendorHandler.VerifyGSTIN)
			// etc.
		}

		// Product routes
		products := v1.Group("/products")
		{
			// TODO: Add product endpoints
			// products.GET("", productHandler.List)
			// products.GET("/:slug", productHandler.GetBySlug)
			// etc.
		}

		// Cart routes
		cart := v1.Group("/cart")
		cart.Use(authMiddleware.RequireAuth())
		{
			// TODO: Add cart endpoints
		}

		// Order routes
		orders := v1.Group("/orders")
		orders.Use(authMiddleware.RequireAuth())
		{
			// TODO: Add order endpoints
		}

		// Admin routes
		admin := v1.Group("/admin")
		admin.Use(authMiddleware.RequireAuth())
		admin.Use(authMiddleware.RequireAdmin())
		{
			// TODO: Add admin endpoints
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in goroutine
	go func() {
		log.Info("Server starting", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Server shutting down...")

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", err)
	}

	log.Info("Server stopped")
}

// healthCheckHandler returns a health check handler
func healthCheckHandler(db *database.Database, cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		health := gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"services":  gin.H{},
		}

		// Check database
		if err := db.Health(); err != nil {
			health["services"].(gin.H)["database"] = "unhealthy"
			health["status"] = "degraded"
		} else {
			health["services"].(gin.H)["database"] = "healthy"
		}

		// Check Redis
		if cache != nil {
			if err := cache.Health(c.Request.Context()); err != nil {
				health["services"].(gin.H)["redis"] = "unhealthy"
				health["status"] = "degraded"
			} else {
				health["services"].(gin.H)["redis"] = "healthy"
			}
		} else {
			health["services"].(gin.H)["redis"] = "not configured"
		}

		if health["status"] == "degraded" {
			c.JSON(http.StatusServiceUnavailable, health)
			return
		}

		c.JSON(http.StatusOK, health)
	}
}

// corsMiddleware configures CORS
func corsMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.CORS.AllowedOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Methods", cfg.CORS.AllowedMethods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", cfg.CORS.AllowedHeaders)
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
