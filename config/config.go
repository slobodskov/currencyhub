// Package config handles application configuration management
// Loads settings from YAML file and environment variables with fallback defaults
package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
)

// Config represents the main application configuration structure
// Contains nested structures for different service configurations
type Config struct {
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     int    `yaml:"port" env:"DB_PORT"`
		User     string `yaml:"user" env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		DBName   string `yaml:"dbname" env:"DB_NAME"`
		SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE"`
	} `yaml:"database"`
	Telegram struct {
		Token string `env:"TELEGRAM_TOKEN"`
	} `yaml:"telegram"`
	Coingecko struct {
		APIKey string `env:"COINGECKO_API_KEY"`
	} `yaml:"coingecko"`
	Server struct {
		Port string `yaml:"port" env:"SERVER_PORT"`
	} `yaml:"server"`
	Logging struct {
		File string `yaml:"file" env:"LOG_FILE"`
	} `yaml:"logging"`
}

// Load initializes and loads application configuration
// Reads from config.yaml file and environment variables
// Returns Config pointer and error if loading fails
func Load() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config

	file, err := os.Open("config/config.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if cfg.Database.Password == "" {
		cfg.Database.Password = "password"
	}

	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
