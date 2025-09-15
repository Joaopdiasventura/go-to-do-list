package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN  string
	Port string
}

func Load() Config {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Config{DSN: dsn, Port: port}
}
