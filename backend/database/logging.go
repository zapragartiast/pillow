package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// LoggingDB wraps sql.DB to log queries
type LoggingDB struct {
	*sql.DB
	logger           *logrus.Logger
	isLoggingEnabled bool
}

// NewLoggingDB creates a new LoggingDB wrapper
func NewLoggingDB(db *sql.DB, logger *logrus.Logger, isLoggingEnabled bool) *LoggingDB {
	return &LoggingDB{
		DB:               db,
		logger:           logger,
		isLoggingEnabled: isLoggingEnabled,
	}
}

// logQuery logs database query details
func (ldb *LoggingDB) logQuery(query string, args []interface{}, start time.Time, err error, rowsAffected int64) {
	if !ldb.isLoggingEnabled || ldb.logger == nil {
		return
	}

	duration := time.Since(start)

	// Sanitize args to avoid logging sensitive data
	sanitizedArgs := make([]interface{}, len(args))
	for i, arg := range args {
		if str, ok := arg.(string); ok {
			// Mask potential passwords or sensitive strings
			if len(str) > 10 && containsSensitiveData(str) {
				sanitizedArgs[i] = "***MASKED***"
			} else {
				sanitizedArgs[i] = arg
			}
		} else {
			sanitizedArgs[i] = arg
		}
	}

	fields := logrus.Fields{
		"level":         "debug",
		"type":          "database_query",
		"query":         query,
		"args":          sanitizedArgs,
		"duration_ms":   duration.Milliseconds(),
		"rows_affected": rowsAffected,
	}

	if err != nil {
		fields["error"] = err.Error()
		ldb.logger.WithFields(fields).Error("Database Query Failed")
	} else {
		ldb.logger.WithFields(fields).Debug("Database Query")
	}
}

// containsSensitiveData checks if string might contain sensitive data
func containsSensitiveData(s string) bool {
	sensitive := []string{"password", "token", "secret", "key", "auth"}
	for _, word := range sensitive {
		if len(s) > len(word) && contains(s, word) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

// Exec logs the query and calls the underlying Exec
func (ldb *LoggingDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := ldb.DB.Exec(query, args...)
	var rowsAffected int64 = 0
	if err == nil && result != nil {
		rowsAffected, _ = result.RowsAffected()
	}
	ldb.logQuery(query, args, start, err, rowsAffected)
	return result, err
}

// ExecContext logs the query and calls the underlying ExecContext
func (ldb *LoggingDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := ldb.DB.ExecContext(ctx, query, args...)
	var rowsAffected int64 = 0
	if err == nil && result != nil {
		rowsAffected, _ = result.RowsAffected()
	}
	ldb.logQuery(query, args, start, err, rowsAffected)
	return result, err
}

// Query logs the query and calls the underlying Query
func (ldb *LoggingDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := ldb.DB.Query(query, args...)
	ldb.logQuery(query, args, start, err, 0) // Rows affected not available for SELECT
	return rows, err
}

// QueryContext logs the query and calls the underlying QueryContext
func (ldb *LoggingDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := ldb.DB.QueryContext(ctx, query, args...)
	ldb.logQuery(query, args, start, err, 0) // Rows affected not available for SELECT
	return rows, err
}

// QueryRow logs the query and calls the underlying QueryRow
func (ldb *LoggingDB) QueryRow(query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := ldb.DB.QueryRow(query, args...)
	// For QueryRow, we can't easily get error/result, so just log the query
	ldb.logQuery(query, args, start, nil, 0)
	return row
}

// QueryRowContext logs the query and calls the underlying QueryRowContext
func (ldb *LoggingDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := ldb.DB.QueryRowContext(ctx, query, args...)
	ldb.logQuery(query, args, start, nil, 0)
	return row
}

// Prepare logs the prepare statement
func (ldb *LoggingDB) Prepare(query string) (*LoggingStmt, error) {
	start := time.Now()
	stmt, err := ldb.DB.Prepare(query)
	loggingStmt := &LoggingStmt{
		Stmt:             stmt,
		query:            query,
		logger:           ldb.logger,
		isLoggingEnabled: ldb.isLoggingEnabled,
	}
	if err == nil {
		ldb.logQuery(fmt.Sprintf("PREPARE: %s", query), nil, start, err, 0)
	}
	return loggingStmt, err
}

// LoggingStmt wraps sql.Stmt to log executions
type LoggingStmt struct {
	*sql.Stmt
	query            string
	logger           *logrus.Logger
	isLoggingEnabled bool
}

// Exec logs the statement execution
func (ls *LoggingStmt) Exec(args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := ls.Stmt.Exec(args...)
	var rowsAffected int64 = 0
	if err == nil && result != nil {
		rowsAffected, _ = result.RowsAffected()
	}
	ls.logStmt(args, start, err, rowsAffected)
	return result, err
}

// Query logs the statement query
func (ls *LoggingStmt) Query(args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := ls.Stmt.Query(args...)
	ls.logStmt(args, start, err, 0)
	return rows, err
}

// QueryRow logs the statement query row
func (ls *LoggingStmt) QueryRow(args ...interface{}) *sql.Row {
	start := time.Now()
	row := ls.Stmt.QueryRow(args...)
	ls.logStmt(args, start, nil, 0)
	return row
}

func (ls *LoggingStmt) logStmt(args []interface{}, start time.Time, err error, rowsAffected int64) {
	if !ls.isLoggingEnabled || ls.logger == nil {
		return
	}

	duration := time.Since(start)

	sanitizedArgs := make([]interface{}, len(args))
	for i, arg := range args {
		if str, ok := arg.(string); ok {
			if len(str) > 10 && containsSensitiveData(str) {
				sanitizedArgs[i] = "***MASKED***"
			} else {
				sanitizedArgs[i] = arg
			}
		} else {
			sanitizedArgs[i] = arg
		}
	}

	fields := logrus.Fields{
		"level":         "debug",
		"type":          "database_statement",
		"statement":     ls.query,
		"args":          sanitizedArgs,
		"duration_ms":   duration.Milliseconds(),
		"rows_affected": rowsAffected,
	}

	if err != nil {
		fields["error"] = err.Error()
		ls.logger.WithFields(fields).Error("Database Statement Failed")
	} else {
		ls.logger.WithFields(fields).Debug("Database Statement")
	}
}
