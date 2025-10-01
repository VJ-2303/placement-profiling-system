package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VJ-2303/placement-profiling-system/internal/auth"
	"github.com/VJ-2303/placement-profiling-system/internal/data"
	"github.com/VJ-2303/placement-profiling-system/internal/models"
	"github.com/joho/godotenv"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
	oauth struct {
		clientID     string
		clientSecret string
		redirectURL  string
	}
	jwt struct {
		secret string
	}
	frontend struct {
		successURL string
	}
}

type application struct {
	config     config
	logger     *log.Logger
	models     models.Models
	msOAuth    *auth.MicrosoftOAuth
	jwtService *auth.JWTService
}

func main() {
	var cfg config

	// Load environment variables only in development
	// In production (Railway), environment variables are already set
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: .env file not found, using environment variables")
		}
	}

	cfg.port = getEnvWithDefault("PORT", "4000")
	cfg.env = getEnvWithDefault("ENV", "development")

	// Load config from environment variables
	cfg.db.dsn = os.Getenv("DB_DSN")
	cfg.oauth.clientID = os.Getenv("CLIENT_ID")
	cfg.oauth.clientSecret = os.Getenv("CLIENT_SECRET")

	// Set redirect URL based on environment
	if cfg.env == "production" {
		cfg.oauth.redirectURL = "https://placement-profiling-system-production.up.railway.app/auth/callback"
	} else {
		cfg.oauth.redirectURL = "http://localhost:4000/auth/callback"
	}

	cfg.jwt.secret = os.Getenv("JWT_SECRET")
	cfg.frontend.successURL = getEnvWithDefault("FRONTEND_SUCCESS_URL", "http://localhost:3000/profile.html")

	// Validate required environment variables
	if cfg.db.dsn == "" {
		log.Fatal("DB_DSN environment variable is required")
	}
	if cfg.oauth.clientID == "" {
		log.Fatal("CLIENT_ID environment variable is required")
	}
	if cfg.oauth.clientSecret == "" {
		log.Fatal("CLIENT_SECRET environment variable is required")
	}
	if cfg.jwt.secret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// Initialize logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Open database connection
	db, err := data.OpenDB(cfg.db.dsn)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Printf("database connection pool established")

	// Initialize application
	app := &application{
		config:     cfg,
		logger:     logger,
		models:     models.NewModels(db),
		msOAuth:    auth.NewMicrosoftOAuth(cfg.oauth.clientID, cfg.oauth.clientSecret, cfg.oauth.redirectURL),
		jwtService: auth.NewJWTService(cfg.jwt.secret),
	}

	// Start server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
