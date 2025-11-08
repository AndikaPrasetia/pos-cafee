package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/AndikaPrasetia/pos-cafee/config"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading: %v", err)
		// Continue execution using environment variables from the system
	}
    cfg := config.LoadDBConfig()

    db, err := sql.Open("postgres", cfg.GetConnectionString())
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Test connection
    err = db.Ping()
    if err != nil {
        log.Fatal("Failed to ping database:", err)
    }

    fmt.Println("✅ Successfully connected to database!")

    // Check users table
    var userCount int
    err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
    if err != nil {
        log.Fatal("Failed to query users:", err)
    }

    fmt.Printf("📊 Found %d users in database\n", userCount)

    // Check menu items
    var menuCount int
    err = db.QueryRow("SELECT COUNT(*) FROM menu_items").Scan(&menuCount)
    if err != nil {
        log.Fatal("Failed to query menu items:", err)
    }

    fmt.Printf("🍽️  Found %d menu items in database\n", menuCount)
}
