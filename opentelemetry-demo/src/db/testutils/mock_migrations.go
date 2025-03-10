package testutils

import (
	"fmt"
	"os"
	"path/filepath"
)

// Test migration SQL content
const (
	initialSchemaMigration = `-- Migration: V1__initial_schema.sql
-- Description: Initial database schema for tests
-- Services: Test Service

-- ==================== UP MIGRATION ====================

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indices
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- ==================== DOWN MIGRATION ====================

-- To roll back this migration, uncomment and execute these statements:
-- DROP INDEX IF EXISTS idx_users_username;
-- DROP TABLE IF EXISTS users;`

	testUserMigration = `-- Migration: V2__add_test_user.sql
-- Description: Adds a test user for testing
-- Services: Test Service

-- ==================== UP MIGRATION ====================

-- Create test user (password: testpassword123)
INSERT INTO users (username, password_hash) 
VALUES ('testuser', '$2a$10$QGfO0JUVuG5R.lQGXSIzd.pBB7WmJjkJ6zf6jE/oyGqhR8tGWRYMG')
ON CONFLICT (username) DO NOTHING;

-- ==================== DOWN MIGRATION ====================

-- To roll back this migration, uncomment and execute this statement:
-- DELETE FROM users WHERE username = 'testuser';`
)

// MockMigrations creates a temporary directory with test migrations
func MockMigrations() (string, func(), error) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "db-migrations-test-*")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp directory: %v", err)
	}

	// Cleanup function to remove the temporary directory
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	// Write test migration files
	if err := writeMigrationFile(tempDir, "V1__initial_schema.sql", initialSchemaMigration); err != nil {
		cleanup()
		return "", nil, err
	}

	if err := writeMigrationFile(tempDir, "V2__add_test_user.sql", testUserMigration); err != nil {
		cleanup()
		return "", nil, err
	}

	return tempDir, cleanup, nil
}

// Helper function to write a migration file
func writeMigrationFile(dir, filename, content string) error {
	path := filepath.Join(dir, filename)
	return os.WriteFile(path, []byte(content), 0644)
}
