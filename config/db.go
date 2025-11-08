package config

import (
    "fmt"
    "os"
    "strconv"
)

type DBConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
}

func LoadDBConfig() *DBConfig {
    port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

    return &DBConfig{
        Host:     os.Getenv("DB_HOST"),
        Port:     port,
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASSWORD"),
        DBName:   os.Getenv("DB_NAME"),
        SSLMode:  os.Getenv("DB_SSLMODE"),
    }
}

func (c *DBConfig) GetConnectionString() string {
    return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
