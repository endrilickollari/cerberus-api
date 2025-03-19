package config

import (
	"os"
	"time"
)

// Config holds all application configuration settings
type Config struct {
	Server ServerConfig
	JWT    JWTConfig
}

// ServerConfig holds HTTP server configurations
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// JWTConfig holds JWT configurations
type JWTConfig struct {
	Secret    []byte
	ExpiresIn time.Duration
}

// NewConfig creates a new configuration from environment variables
func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  time.Second * 15,
			WriteTimeout: time.Second * 15,
			IdleTimeout:  time.Second * 60,
		},
		JWT: JWTConfig{
			Secret:    []byte(getEnv("JWT_SECRET", "your_secret_key_please_change_in_production")),
			ExpiresIn: time.Hour * 24,
		},
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
