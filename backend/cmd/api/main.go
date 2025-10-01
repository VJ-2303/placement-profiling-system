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

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg.port = getEnvWithDefault("PORT", "4000")
	cfg.env = getEnvWithDefault("ENV", "development")

	// Load config from environment variables
	cfg.db.dsn = os.Getenv("DB_DSN")
	cfg.oauth.clientID = os.Getenv("CLIENT_ID")
	cfg.oauth.clientSecret = os.Getenv("CLIENT_SECRET")
	cfg.oauth.redirectURL = "https://placement-profiling-system-production.up.railway.app/auth/callback"
	cfg.jwt.secret = os.Getenv("JWT_SECRET")
	cfg.frontend.successURL = getEnvWithDefault("FRONTEND_SUCCESS_URL", "")

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
