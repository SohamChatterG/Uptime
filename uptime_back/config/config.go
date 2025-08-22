package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	MongoURI      string
	DBName        string
	JWTSecret     string
	CheckInterval time.Duration
	EmailUser     string
	EmailPass     string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	intervalStr := os.Getenv("CHECK_INTERVAL_SECONDS")
	interval, err := strconv.Atoi(intervalStr)
	if err != nil || interval <= 0 {
		interval = 60
	}

	emailUser := os.Getenv("EMAIL_USER")
	emailPass := os.Getenv("EMAIL_PASS")
	if emailUser == "" || emailPass == "" {
		log.Println("WARNING: EMAIL_USER or EMAIL_PASS not set. Email notifications will be disabled.")
	}

	return &Config{
		Port:          port,
		MongoURI:      os.Getenv("MONGODB_URI"),
		DBName:        os.Getenv("DB_NAME"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		CheckInterval: time.Duration(interval) * time.Second,
		EmailUser:     emailUser,
		EmailPass:     emailPass,
	}
}
