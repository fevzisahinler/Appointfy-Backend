package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	PGHost       string
	PGUser       string
	PGDBName     string
	PGPassword   string
	PGPort       string
	JwtSecretKey string
}

var cfg *Config

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg = &Config{
		PGHost:       os.Getenv("POSTGRES_HOST"),
		PGUser:       os.Getenv("POSTGRES_USER"),
		PGDBName:     os.Getenv("POSTGRES_DB"),
		PGPassword:   os.Getenv("POSTGRES_PASSWORD"),
		PGPort:       os.Getenv("POSTGRES_PORT"),
		JwtSecretKey: os.Getenv("JwtSecretKey"),
	}
	return cfg, nil
}
