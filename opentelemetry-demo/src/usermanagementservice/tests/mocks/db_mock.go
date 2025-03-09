package mocks

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

// MockDB returns a new sql.DB mock and sqlmock.Sqlmock to configure it
func MockDB() (*sql.DB, sqlmock.Sqlmock, error) {
	return sqlmock.New()
}
