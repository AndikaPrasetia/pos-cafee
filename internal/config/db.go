package config

import (
	"database/sql"
	"log"

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
