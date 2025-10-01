package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	jwt struct {
		secret string
	}
	microsoft struct {
		clientID     string
		clientSecret string
		redirectURL  string
	}
}

func loadConfig() config {
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found")
	}

	var cfg config

	cfg.port = 4001
	cfg.env = "devolopment"
	cfg.db.dsn = os.Getenv("DB_DSN")

	cfg.jwt.secret = os.Getenv("JWT_SECRET")

	cfg.microsoft.clientID = os.Getenv("CLIENT_ID")
	cfg.microsoft.clientSecret = os.Getenv("CLIENT_SECRET")
	cfg.microsoft.redirectURL = "http://localhost:3000/auth/callback"

	return cfg
}
