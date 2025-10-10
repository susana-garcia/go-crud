package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds all configuration for our application
type Config struct {
	Server
	Database
	LogLevel string
	Debug    bool
}

type Server struct {
	Port string
	Host string
}

type Database struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// Load reads configuration from environment variables
func Load() *Config {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Fatal("unable to load env vars: %w", err)
	}

	config := &Config{
		Server: Server{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "localhost"),
		},
		Database: Database{
			Port:     getEnv("DB_PORT", "5432"),
			Host:     getEnv("DB_HOST", "localhost"),
			Name:     getEnv("DB_NAME", "go-crud"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Debug:    getEnvBool("DEBUG", false),
	}

	log.Printf("configuration loaded: port=%s, host=%s, log_level=%s, debug=%t",
		config.Server.Port, config.Server.Host, config.LogLevel, config.Debug)

	return config
}

func OpenConnection(cfg Database) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)

	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("unable to connect to the database: %w", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("unable to return database: %w", err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("database connected")
	return database
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
