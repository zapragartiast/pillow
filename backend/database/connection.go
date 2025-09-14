package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var DB *LoggingDB

func Connect(dbURL string) *LoggingDB {
	var err error
	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	// For now, return regular sql.DB wrapped, but we'll update when logger is available
	DB = &LoggingDB{DB: sqlDB, isLoggingEnabled: false}
	return DB
}

// ConnectWithLogging creates a logging database connection
func ConnectWithLogging(dbURL string, logger *logrus.Logger, isLoggingEnabled bool) *LoggingDB {
	var err error
	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return NewLoggingDB(sqlDB, logger, isLoggingEnabled)
}
