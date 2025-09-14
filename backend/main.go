package main

import (
	"log"
	"net/http"
	"os"
	"pillow/audit"
	"pillow/database"
	"pillow/routes"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
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

	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z07:00",
	})

	// Check environment for logging
	env := os.Getenv("ENV")
	isLoggingEnabled := env == "development"
	if isLoggingEnabled {
		logger.SetLevel(logrus.DebugLevel)
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetLevel(logrus.ErrorLevel)
		logger.SetOutput(os.Stderr)
	}

	db := database.ConnectWithLogging(dbURL, logger, isLoggingEnabled)
	defer db.Close()

	// Start the async audit queue (buffer size 100). This ensures audit events are enqueued
	// by middleware and persisted by the background worker.
	// Extract the underlying *sql.DB from LoggingDB for audit queue
	audit.StartAuditQueue(db.DB, 100)
	defer audit.StopAuditQueue()

	r := routes.SetupRoutes(db, logger, isLoggingEnabled)

	log.Printf("Backend running on :%s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, r); err != nil {
		log.Fatal(err)
	}
}
