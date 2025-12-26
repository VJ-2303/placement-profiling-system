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
		tenantID     string
		redirectURL  string
	}
	jwt struct {
		secret string
	}
	frontend struct {
		url string
	}
	allowedDomain string // e.g., kct.ac.in
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

	cfg.port = getEnvWithDefault("PORT", "4000")
	cfg.env = getEnvWithDefault("ENV", "development")
	cfg.allowedDomain = getEnvWithDefault("ALLOWED_DOMAIN", "kct.ac.in")

	// Database - support both DATABASE_URL (Neon/Railway style) and DB_DSN
	cfg.db.dsn = os.Getenv("DATABASE_URL")
	if cfg.db.dsn == "" {
		cfg.db.dsn = os.Getenv("DB_DSN")
	}

	// Microsoft OAuth - support both naming conventions
	cfg.oauth.clientID = os.Getenv("MICROSOFT_CLIENT_ID")
	if cfg.oauth.clientID == "" {
		cfg.oauth.clientID = os.Getenv("CLIENT_ID")
	}

	cfg.oauth.clientSecret = os.Getenv("MICROSOFT_CLIENT_SECRET")
	if cfg.oauth.clientSecret == "" {
		cfg.oauth.clientSecret = os.Getenv("CLIENT_SECRET")
	}

	cfg.oauth.tenantID = os.Getenv("MICROSOFT_TENANT_ID")

	cfg.oauth.redirectURL = os.Getenv("MICROSOFT_REDIRECT_URL")
	if cfg.oauth.redirectURL == "" {
		cfg.oauth.redirectURL = getEnvWithDefault("OAUTH_REDIRECT_URL", "http://localhost:4000/auth/callback")
	}

	cfg.jwt.secret = os.Getenv("JWT_SECRET")
	cfg.frontend.url = getEnvWithDefault("FRONTEND_URL", "http://localhost:5500")

	// Validate required env vars
	if cfg.db.dsn == "" {
		log.Fatal("DATABASE_URL or DB_DSN environment variable is required")
	}
	if cfg.oauth.clientID == "" {
		log.Fatal("MICROSOFT_CLIENT_ID or CLIENT_ID environment variable is required")
	}
	if cfg.oauth.clientSecret == "" {
		log.Fatal("MICROSOFT_CLIENT_SECRET or CLIENT_SECRET environment variable is required")
	}
	if cfg.jwt.secret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// Initialize logger
	logger := log.New(os.Stdout, "[PPS] ", log.Ldate|log.Ltime|log.Lshortfile)

	logger.Println("=== Placement Profiling System ===")
	logger.Printf("Environment: %s", cfg.env)
	logger.Printf("Port: %s", cfg.port)
	logger.Printf("Allowed Domain: %s", cfg.allowedDomain)
	logger.Printf("Frontend URL: %s", cfg.frontend.url)
	logger.Printf("OAuth Redirect: %s", cfg.oauth.redirectURL)

	// Open database connection
	db, err := data.OpenDB(cfg.db.dsn)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Println("Database connection established")

	// Initialize application struct
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
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	logger.Printf("Starting %s server on port %s", cfg.env, cfg.port)
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
