package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Database struct {
	DSN string
}

type Config struct {
	Port     string
	Database Database
	SECRET   string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, using default config")
	}

	return &Config{
		Port: os.Getenv("SERVER_PORT"),
		Database: Database{
			DSN: os.Getenv("DSN"),
		},
		SECRET: os.Getenv("SECRET"),
	}
}
