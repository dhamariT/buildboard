package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	Environment string // "development" or "production"
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string

	Port        string
	GinMode     string
	FrontendURL string
	BackendURL  string

	// Email (SMTP)
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// IsProduction returns true if running in production environment.
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment returns true if running in development environment.
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// Load reads configuration from environment variables and .env files.
// It attempts to load .env files from common locations, falling back to
// environment variables for production deployments.
func Load() *Config {
	// Try loading .env from common locations
	loaded := tryLoadEnv(
		".env",
		"../.env",
		"../../.env",
		"backend/.env",
		"../backend/.env",
		"../../backend/.env",
	)
	if loaded != "" {
		slog.Info("Environment file loaded", "path", loaded)
	}

	// Detect environment - defaults to development if not set
	env := getEnv("ENVIRONMENT", "development")
	if env != "production" && env != "development" {
		env = "development"
	}

	cfg := &Config{
		Environment: env,
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "buildboard_db"),
		DBSSLMode:   getEnv("DB_SSL_MODE", "disable"),
		Port:        getEnv("PORT", "8080"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		BackendURL:  getEnv("BACKEND_URL", "http://localhost:8080"),
		// Email (SMTP)
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnvInt("SMTP_PORT", 587),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		FromEmail:    getEnv("FROM_EMAIL", ""),
		FromName:     getEnv("FROM_NAME", "BuildBoard"),
	}

	// Warn about default/empty sensitive values
	if cfg.DBPassword == "postgres" || cfg.DBPassword == "" {
		slog.Warn("Using default or empty DB password; set DB_PASSWORD in env for production")
	}

	// Validate required variables for production
	if cfg.IsProduction() {
		if err := validateProduction(cfg); err != nil {
			slog.Error("Production environment validation failed", "error", err)
			os.Exit(1)
		}
		slog.Info("Production environment validation passed")
	}

	return cfg
}

// validateProduction ensures all required environment variables are set for production.
func validateProduction(cfg *Config) error {
	var missing []string

	// Database - required for production
	if cfg.DBHost == "" {
		missing = append(missing, "DB_HOST")
	}
	if cfg.DBPort == "" {
		missing = append(missing, "DB_PORT")
	}
	if cfg.DBUser == "" {
		missing = append(missing, "DB_USER")
	}
	if cfg.DBPassword == "" || cfg.DBPassword == "postgres" {
		missing = append(missing, "DB_PASSWORD")
	}
	if cfg.DBName == "" {
		missing = append(missing, "DB_NAME")
	}

	// Server - required for production
	if cfg.Port == "" {
		missing = append(missing, "PORT")
	}
	if cfg.FrontendURL == "" {
		missing = append(missing, "FRONTEND_URL")
	}
	if cfg.BackendURL == "" {
		missing = append(missing, "BACKEND_URL")
	}

	// Email (SMTP) - required for production
	if cfg.SMTPHost == "" {
		missing = append(missing, "SMTP_HOST")
	}
	if cfg.SMTPPort == 0 {
		missing = append(missing, "SMTP_PORT")
	}
	if cfg.SMTPUsername == "" {
		missing = append(missing, "SMTP_USERNAME")
	}
	if cfg.SMTPPassword == "" {
		missing = append(missing, "SMTP_PASSWORD")
	}
	if cfg.FromEmail == "" {
		missing = append(missing, "FROM_EMAIL")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables for production: %s", strings.Join(missing, ", "))
	}

	return nil
}

func tryLoadEnv(paths ...string) string {
	for _, p := range paths {
		if err := godotenv.Load(p); err == nil {
			return p
		}
	}
	return ""
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
