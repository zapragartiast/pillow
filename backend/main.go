package main

import (
	"log"
	"net/http"
	"os"
	"pillow/audit"
	"pillow/database"
	"pillow/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080" // Default port
	}

	db := database.Connect(dbURL)
	defer db.Close()

	// Start the async audit queue (buffer size 100). This ensures audit events are enqueued
	// by middleware and persisted by the background worker.
	audit.StartAuditQueue(db, 100)
	defer audit.StopAuditQueue()

	r := routes.SetupRoutes(db)

	log.Printf("Backend running on :%s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, r); err != nil {
		log.Fatal(err)
	}
}
