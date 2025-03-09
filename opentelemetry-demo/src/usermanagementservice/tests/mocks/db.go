package mocks

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
)

// This file will contain mock implementations for database interactions
// It should be filled with actual mock implementations when writing tests 

// NewDBMock creates a new mock database and sqlmock for testing
func NewDBMock() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create sqlmock: %v", err)
	}

	return db, mock, nil
}

// ExpectUsernameCheck sets up the mock to expect a username check query
func ExpectUsernameCheck(mock sqlmock.Sqlmock, username string, exists bool) {
	rows := sqlmock.NewRows([]string{"exists"}).AddRow(exists)
	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE username = \\$1\\)").
		WithArgs(username).
		WillReturnRows(rows)
}

// ExpectInsertUser sets up the mock to expect a user insertion query
func ExpectInsertUser(mock sqlmock.Sqlmock, username, hashedPassword string, userId int32) {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(userId)
	mock.ExpectQuery("INSERT INTO users \\(username, password_hash\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
		WithArgs(username, sqlmock.AnyArg()).
		WillReturnRows(rows)
}

// ExpectGetUser sets up the mock to expect a user retrieval query
func ExpectGetUser(mock sqlmock.Sqlmock, username string, userID int32, storedPasswordHash string) {
	rows := sqlmock.NewRows([]string{"id", "username", "password_hash"}).
		AddRow(userID, username, storedPasswordHash)
	mock.ExpectQuery("SELECT id, username, password_hash FROM users WHERE username = \\$1").
		WithArgs(username).
		WillReturnRows(rows)
}

// ExpectGetUserNotFound sets up the mock to expect a user retrieval query that returns no rows
func ExpectGetUserNotFound(mock sqlmock.Sqlmock, username string) {
	mock.ExpectQuery("SELECT id, username, password_hash FROM users WHERE username = \\$1").
		WithArgs(username).
		WillReturnError(sql.ErrNoRows)
}

// ExpectError sets up the mock to expect any query and return an error
func ExpectError(mock sqlmock.Sqlmock, queryPattern string, args ...driver.Value) {
	mock.ExpectQuery(queryPattern).
		WithArgs(args...).
		WillReturnError(fmt.Errorf("database error"))
} 