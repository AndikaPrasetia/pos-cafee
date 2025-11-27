package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	URL      string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// AppConfig holds application configuration
type AppConfig struct {
	Environment string
	Port        string
	JWTSecret   string
	JWTExpiry   string
	LogLevel    string
	DB          DBConfig
	Redis       RedisConfig
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *AppConfig {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Set default values
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENV", "development")
	viper.SetDefault("JWT_EXPIRY", "24h")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSL_MODE", "disable")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_DB", "0")

	// Read configuration from environment variables
	viper.AutomaticEnv()

	// Create config struct
	config := &AppConfig{
		Environment: getEnv("APP_ENV", "development"),
		Port:        getEnv("APP_PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "default-secret-key-for-development-do-not-use-in-production"),
		JWTExpiry:   getEnv("JWT_EXPIRY", "24h"),
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "username"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "pos_cafe"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			URL:      getEnv("DATABASE_URL", ""),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
	}

	// If DATABASE_URL is not set, construct it from individual components
	if config.DB.URL == "" {
		config.DB.URL = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
			config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.DBName, config.DB.SSLMode)
	}

	return config
}

// getEnv retrieves environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves environment variable as integer or returns default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
