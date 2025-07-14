package main

import (
	"log"
	"os"

	"inventory_system/database"
	"inventory_system/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database
	database.InitDatabase()

	// Setup routes
	router := routes.SetupRoutes()

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get host from environment (default to 0.0.0.0 for all interfaces)
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	address := host + ":" + port
	log.Printf("Server starting on %s", address)
	log.Fatal(router.Run(address))
}
