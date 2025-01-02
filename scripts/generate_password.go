package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/muety/wakapi/utils"
)

func main() {
	// Define command line flags
	password := flag.String("password", "", "The plain password to hash")
	salt := flag.String("salt", "", "The password salt (must match your server's configuration)")

	// Parse flags
	flag.Parse()

	// Validate inputs
	if *password == "" || *salt == "" {
		fmt.Println("Usage: generate_password -password <password> -salt <salt>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Generate hash
	hash, err := utils.HashPassword(*password, *salt)
	if err != nil {
		log.Fatalf("Failed to generate password hash: %v", err)
	}

	// Output the hash
	fmt.Printf("Password Hash: %s\n", hash)
}
