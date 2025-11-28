package config

import (
	"context"
	"database/sql"
	"log"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

// ConnectDB creates and returns a database connection
func ConnectDB(config *AppConfig) *sql.DB {
	db, err := sql.Open("postgres", config.DB.URL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	log.Println("Successfully connected to database")
	return db
}

// CloseDB closes the database connection
func CloseDB(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}
}

// RedisClient creates and returns a Redis client
func RedisClient(config *AppConfig) *redis.Client {
	// If REDIS_URL is provided, use it to create the client
	if config.Redis.URL != "" {
		opts, err := redis.ParseURL(config.Redis.URL)
		if err != nil {
			log.Fatal("Failed to parse Redis URL:", err)
		}

		rdb := redis.NewClient(opts)

		// Test the connection
		ctx := context.Background()
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Fatal("Failed to connect to Redis:", err)
		}

		log.Println("Successfully connected to Redis using REDIS_URL")
		return rdb
	}

	// Fallback to individual config values if REDIS_URL is not provided
	// This is mainly for backward compatibility during the transition
	log.Println("REDIS_URL not provided, using individual Redis configuration (deprecated)")

	// For backward compatibility, we'll assume local Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Test the connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
			log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Successfully connected to Redis")
	return rdb
}

// CloseRedis closes the Redis client connection
func CloseRedis(rdb *redis.Client) {
	if rdb != nil {
		if err := rdb.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
	}
}
