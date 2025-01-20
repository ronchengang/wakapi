package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql" // or your database driver
	"github.com/muety/wakapi/utils"
)

type User struct {
	ID                string
	email             string
	has_data          bool
	created_at        string
	last_logged_in_at string
}

func main() {
	// Define command line flags
	password := flag.String("password", "", "The plain password to hash")
	salt := flag.String("salt", "", "The password salt (must match your server's configuration)")
	email := flag.String("email", "", "The email to search for users")
	dbConn := flag.String("db", "user:password@tcp(localhost:3306)/wakapi", "Database connection string")

	// Parse flags
	flag.Parse()

	// Validate inputs
	if *password == "" || *email == "" {
		fmt.Println("Usage: generate_password -password <password> -salt <salt> -email <email> [-db <connection_string>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Connect to database
	db, err := sql.Open("mysql", *dbConn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Find users with matching email
	rows, err := db.Query("SELECT id, email,has_data,created_at,last_logged_in_at FROM users WHERE email = ?", *email)
	if err != nil {
		log.Fatalf("Failed to query users: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.email, &user.has_data, &user.created_at, &user.last_logged_in_at); err != nil {
			log.Fatalf("Failed to scan user row: %v", err)
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		log.Fatalf("No users found with email: %s", *email)
	}

	// Display found users
	fmt.Println("\nFound users:")
	for i, user := range users {
		fmt.Printf("[%d] ID: %s, email: %s, has_data: %d, created_at: %s, last_logged_in_at: %s\n", i+1, user.ID, user.email, user.has_data, user.created_at, user.last_logged_in_at)
	}

	// Get user selection
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEnter the number of the user to update (or 'q' to quit): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" {
		fmt.Println("Operation cancelled")
		return
	}

	selection, err := strconv.Atoi(input)
	if err != nil || selection < 1 || selection > len(users) {
		log.Fatalf("Invalid selection. Please choose a number between 1 and %d", len(users))
	}

	selectedUser := users[selection-1]

	// Generate hash
	hash, err := utils.HashPassword(*password, *salt)
	if err != nil {
		log.Fatalf("Failed to generate password hash: %v", err)
	}

	// Update password in database
	result, err := db.Exec("UPDATE users SET password = ? WHERE id = ?", hash, selectedUser.ID)
	if err != nil {
		log.Fatalf("Failed to update password: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		log.Fatalf("Failed to update password for user: %s", selectedUser.email)
	}

	// Output results
	fmt.Printf("Successfully updated password for user:\n")
	fmt.Printf("ID: %s, email: %s\n", selectedUser.ID, selectedUser.email)
}
