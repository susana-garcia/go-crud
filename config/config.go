package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	Port     string
	Host     string
	LogLevel string
	Debug    bool
}

// Load reads configuration from environment variables
func Load() *Config {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	godotenv.Load(".env." + env)

	config := &Config{
		Port:     getEnv("PORT", "8080"),
		Host:     getEnv("HOST", "localhost"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Debug:    getEnvBool("DEBUG", false),
	}

	log.Printf("configuration loaded: port=%s, host=%s, log_level=%s, debug=%t",
		config.Port, config.Host, config.LogLevel, config.Debug)

	return config
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool gets a boolean environment variable or returns a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
