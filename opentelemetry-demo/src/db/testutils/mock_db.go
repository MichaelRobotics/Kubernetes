package testutils

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

// MockDB creates a mock database for testing
func MockDB() (*sql.DB, sqlmock.Sqlmock, error) {
	return sqlmock.New()
}
