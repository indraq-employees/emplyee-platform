package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	MongoURI    string
	MongoDB     string
	JWTSecret   string
	FrontendURL string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:     getEnv("MONGO_DB", "employee_platform"),
		JWTSecret:   getEnv("JWT_SECRET", "change_this_super_secret"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}