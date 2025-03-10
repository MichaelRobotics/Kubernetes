// Package postgres provides database connection utilities for PostgreSQL
package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Connection represents a PostgreSQL database connection
type Connection struct {
	DB *sql.DB
}

// NewConnection creates a new PostgreSQL database connection
func NewConnection(connString string) (*Connection, error) {
	if connString == "" {
		return nil, fmt.Errorf("database connection string is required")
	}

	// Open connection to database
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection works
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Connection{DB: db}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

// GetConnectionFromEnv creates a database connection from environment variables
func GetConnectionFromEnv(envVar string) (*Connection, error) {
	connString := os.Getenv(envVar)
	if connString == "" {
		return nil, fmt.Errorf("%s environment variable not set", envVar)
	}

	return NewConnection(connString)
}

// GetConnection creates a direct sql.DB connection to the database
func GetConnection(connString string) (*sql.DB, error) {
	conn, err := NewConnection(connString)
	if err != nil {
		return nil, err
	}
	return conn.DB, nil
}
