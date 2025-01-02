package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/muety/wakapi/config"
	"github.com/muety/wakapi/models"
	"github.com/muety/wakapi/utils"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Define command line flags
	password := flag.String("password", "", "The plain password to hash")
	salt := flag.String("salt", "", "The password salt (must match your server's configuration)")
	userID := flag.String("user", "", "The user ID to update")

	// Parse flags
	flag.Parse()

	// Validate inputs
	if *password == "" || *salt == "" || *userID == "" {
		fmt.Println("Usage: generate_password -password <password> -salt <salt> -user <user_id>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Generate hash
	hash, err := utils.HashPassword(*password, *salt)
	if err != nil {
		log.Fatalf("Failed to generate password hash: %v", err)
	}

	// Load config
	cfg := config.Get()
	if cfg == nil {
		log.Fatal("Failed to load configuration")
	}

	// Initialize database connection
	var db *gorm.DB
	switch cfg.Database.Dialect {
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(cfg.Database.Connection), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(cfg.Database.Connection), &gorm.Config{})
	case "postgres":
		db, err = gorm.Open(postgres.Open(cfg.Database.Connection), &gorm.Config{})
	default:
		log.Fatalf("Unsupported database dialect: %s", cfg.Database.Dialect)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Update user password in database
	result := db.Model(&models.User{}).
		Where("id = ?", *userID).
		Update("password", hash)

	if result.Error != nil {
		log.Fatalf("Failed to update password in database: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Fatalf("No user found with ID: %s", *userID)
	}

	fmt.Printf("Successfully updated password for user %s\n", *userID)
	fmt.Printf("Password Hash: %s\n", hash)
}
