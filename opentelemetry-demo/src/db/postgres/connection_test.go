package postgres

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewConnection(t *testing.T) {
	// Save original function
	originalSqlOpen := sqlOpen
	// Restore original function after test
	defer func() {
		sqlOpen = originalSqlOpen
	}()

	t.Run("SuccessfulConnection", func(t *testing.T) {
		// Create a mock DB using sqlmock directly
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up the mock to expect a ping
		mock.ExpectPing()

		// Override sqlOpen to return our mock DB
		sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
			return db, nil
		}

		// Test successful connection
		conn, err := NewConnection("postgres://user:pass@localhost:5432/testdb")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if conn == nil {
			t.Error("Expected connection, got nil")
		}
		if conn.DB != db {
			t.Error("Connection DB does not match expected DB")
		}

		// Verify that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %s", err)
		}
	})

	t.Run("EmptyConnectionString", func(t *testing.T) {
		// Test connection with invalid connection string
		conn, err := NewConnection("")

		if err == nil {
			t.Error("Expected error for empty connection string, got nil")
		}
		if conn != nil {
			t.Errorf("Expected nil connection, got: %v", conn)
		}
	})

	t.Run("UnreachableDatabase", func(t *testing.T) {
		// Create a mock DB using sqlmock directly
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up the mock to expect and fail a ping
		pingErr := errors.New("ping failed: database unreachable")
		mock.ExpectPing().WillReturnError(pingErr)

		// Override sqlOpen to return our mock DB
		sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
			return db, nil
		}

		// Test connection with unreachable database
		conn, err := NewConnection("postgres://user:pass@nonexistent:5432/testdb")

		if err == nil {
			t.Error("Expected error for unreachable database, got nil")
		}
		if conn != nil {
			t.Errorf("Expected nil connection, got: %v", conn)
		}

		// Verify that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %s", err)
		}
	})

	t.Run("SqlOpenError", func(t *testing.T) {
		// Override sqlOpen to return an error
		sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
			return nil, errors.New("failed to open database")
		}

		// Test connection with sql.Open error
		conn, err := NewConnection("postgres://invalid")

		if err == nil {
			t.Error("Expected error from sql.Open, got nil")
		}
		if conn != nil {
			t.Errorf("Expected nil connection, got: %v", conn)
		}
	})
}

func TestGetConnectionFromEnv(t *testing.T) {
	// Save original function
	originalSqlOpen := sqlOpen
	// Restore original function after test
	defer func() {
		sqlOpen = originalSqlOpen
	}()

	// Environment variable name to use for tests
	testEnvVar := "TEST_DB_CONN"

	t.Run("EnvVarSet", func(t *testing.T) {
		// Set environment variable
		connString := "postgres://user:pass@localhost:5432/testdb"
		os.Setenv(testEnvVar, connString)
		defer os.Unsetenv(testEnvVar)

		// Create a mock DB using sqlmock directly
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up the mock to expect a ping
		mock.ExpectPing()

		// Override sqlOpen to return our mock DB
		sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
			return db, nil
		}

		// Test with environment variable set
		conn, err := GetConnectionFromEnv(testEnvVar)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if conn == nil {
			t.Error("Expected connection, got nil")
		}

		// Verify that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %s", err)
		}
	})

	t.Run("EnvVarNotSet", func(t *testing.T) {
		// Make sure environment variable is not set
		envVarName := "NONEXISTENT_DB_CONN"
		os.Unsetenv(envVarName)

		// Test with environment variable not set
		conn, err := GetConnectionFromEnv(envVarName)

		if err == nil {
			t.Error("Expected error for missing environment variable, got nil")
		}
		if conn != nil {
			t.Errorf("Expected nil connection, got: %v", conn)
		}
	})

	t.Run("EnvVarWithInvalidConnString", func(t *testing.T) {
		// Set environment variable with invalid connection string
		invalidConnString := "invalid_connection_string"
		os.Setenv(testEnvVar, invalidConnString)
		defer os.Unsetenv(testEnvVar)

		// Override sqlOpen to return an error for the invalid string
		sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
			return nil, errors.New("invalid connection string")
		}

		// Test with invalid connection string in environment variable
		conn, err := GetConnectionFromEnv(testEnvVar)

		if err == nil {
			t.Error("Expected error for invalid connection string, got nil")
		}
		if conn != nil {
			t.Errorf("Expected nil connection, got: %v", conn)
		}
	})
}

func TestConnectionClose(t *testing.T) {
	t.Run("CloseOpenConnection", func(t *testing.T) {
		// Create a mock connection
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}

		// Create a Connection with the mock DB
		conn := &Connection{DB: db}

		// Set expectations for Close
		mock.ExpectClose()

		// Test closing an open connection
		err = conn.Close()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Verify that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %s", err)
		}
	})

	t.Run("CloseNilConnection", func(t *testing.T) {
		// Create a Connection with nil DB
		conn := &Connection{DB: nil}

		// Test closing a nil connection
		err := conn.Close()
		if err != nil {
			t.Errorf("Expected no error closing nil connection, got: %v", err)
		}
	})

	t.Run("CloseWithError", func(t *testing.T) {
		// Create a mock connection
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}

		// Create a Connection with the mock DB
		conn := &Connection{DB: db}

		// Set expectations for Close with error
		closeErr := errors.New("close error")
		mock.ExpectClose().WillReturnError(closeErr)

		// Test closing a connection with error
		err = conn.Close()
		if err != closeErr {
			t.Errorf("Expected close error, got: %v", err)
		}

		// Verify that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %s", err)
		}
	})
}
