package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/buildboard/backend/config"
	"github.com/buildboard/backend/internal/controllers"
	"github.com/buildboard/backend/internal/models"
	"github.com/buildboard/backend/internal/services"
	"github.com/buildboard/backend/pkg/db"
)

func main() {
	cfg := config.Load()

	// Print environment information
	slog.Info("BuildBoard Backend Starting")
	slog.Info("Environment configuration", "environment", cfg.Environment)
	if cfg.IsProduction() {
		slog.Info("Running in PRODUCTION mode")
	} else {
		slog.Info("Running in DEVELOPMENT mode")
	}

	gin.SetMode(cfg.GinMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// CORS configuration
	corsCfg := cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsCfg))

	// In development, create database if it doesn't exist
	if cfg.IsDevelopment() {
		if err := db.CreateDatabaseIfNotExists(cfg); err != nil {
			slog.Warn("Failed to create database in development", "error", err)
		}
	}

	// Database connection
	slog.Info("Attempting database connection")
	database, err := db.Connect(cfg)
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		slog.Info("Continuing without database for health checks")
		database = nil
	} else {
		// Ensure DB is actually reachable by executing a simple query
		var one int
		if err := database.Raw("SELECT 1").Scan(&one).Error; err != nil {
			slog.Error("Database ping failed", "error", err)
			slog.Info("Continuing without database for health checks")
			database = nil
		} else {
			slog.Info("Database connection established")
			// Run migrations
			if err := database.AutoMigrate(&models.EarlyStartUser{}); err != nil {
				slog.Error("Database migration failed", "error", err)
				if cfg.IsDevelopment() {
					slog.Error("Migration failure in development - exiting")
					os.Exit(1)
				} else {
					slog.Info("Continuing without migrations in production")
				}
			} else {
				slog.Info("Database migrations completed")
			}
		}
	}

	// Initialize email service
	emailService, err := services.NewEmailService(cfg)
	if err != nil {
		slog.Error("Failed to initialize email service", "error", err)
		emailService = nil
	}

	// Initialize controllers
	earlyStartCtl := controllers.NewEarlyStartController(database, emailService)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		status := "healthy"
		dbStatus := "disconnected"

		if database != nil {
			var one int
			if err := database.Raw("SELECT 1").Scan(&one).Error; err == nil {
				dbStatus = "connected"
			} else {
				dbStatus = "error"
				status = "degraded"
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   status,
			"database": dbStatus,
			"version":  "1.0.0",
		})
	})

	// Public alias for health under /api
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// Early start signup endpoints (public)
	r.POST("/early-start/signup", earlyStartCtl.Signup)
	r.POST("/early-start/verify", earlyStartCtl.VerifyOTP)
	r.GET("/early-start/count", earlyStartCtl.Count)

	// Email engagement tracking endpoint
	engagementGroup := r.Group("/api/e")
	{
		engagementGroup.GET("/:filename", earlyStartCtl.TrackEngagement)
	}

	// Admin endpoints (would typically require authentication)
	// For now, these are public but should be protected in production
	adminGroup := r.Group("/admin")
	{
		adminGroup.GET("/early-start", earlyStartCtl.List)
	}

	port := cfg.Port
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	addr := ":" + port

	// Enhanced logging
	slog.Info("Server configuration",
		"port_env", os.Getenv("PORT"),
		"config_port", cfg.Port,
		"final_port", port,
		"listen_address", addr,
		"gin_mode", cfg.GinMode,
		"frontend_url", cfg.FrontendURL,
	)

	backendURL := "http://localhost" + addr
	slog.Info("Server starting", "url", backendURL)

	// Start server
	if err := http.ListenAndServe(addr, r); err != nil {
		slog.Error("Server failed to start", "error", err, "address", addr)
		os.Exit(1)
	}
}
