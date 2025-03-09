package mocks

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMockDB(t *testing.T) {
	// Create a mock database
	db, mock, err := MockDB()

	// Assert that creation was successful
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, mock)

	// Example of how to use mock
	mock.ExpectExec("INSERT INTO users").
		WithArgs("testuser", "password").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock verification would happen in real tests
	// This test just demonstrates that mock creation works
	db.Close()
}
