package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()

	// Set JSON formatter
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z07:00",
	})

	// Check environment
	env := os.Getenv("ENV")
	if env == "development" {
		Logger.SetLevel(logrus.DebugLevel)
		Logger.SetOutput(os.Stdout)
	} else {
		// In production, disable logging or set to error only
		Logger.SetLevel(logrus.ErrorLevel)
		// Could write to file, but for now, disable
		Logger.SetOutput(os.Stderr)
	}
}

// IsLoggingEnabled checks if logging is enabled based on ENV
func IsLoggingEnabled() bool {
	env := os.Getenv("ENV")
	return env == "development"
}
