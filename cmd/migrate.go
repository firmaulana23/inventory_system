package main

import (
	"fmt"
	"inventory_system/database"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run cmd/migrate.go create <migration_name>")
		fmt.Println("  go run cmd/migrate.go run")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a migration name")
			fmt.Println("Example: go run cmd/migrate.go create add_new_column")
			os.Exit(1)
		}
		migrationName := os.Args[2]
		err := database.CreateMigration(migrationName)
		if err != nil {
			log.Fatal("Failed to create migration:", err)
		}

	case "run":
		// Initialize database connection
		database.InitDatabase()
		
		// Run migrations
		err := database.RunMigrations()
		if err != nil {
			log.Fatal("Failed to run migrations:", err)
		}
		fmt.Println("Migrations completed successfully")

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: create, run")
		os.Exit(1)
	}
}