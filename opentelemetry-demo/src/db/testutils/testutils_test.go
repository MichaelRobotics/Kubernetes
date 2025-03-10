package testutils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMockDB(t *testing.T) {
	db, mock, err := MockDB()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	if db == nil {
		t.Error("MockDB returned nil database")
	}

	if mock == nil {
		t.Error("MockDB returned nil mock")
	}
}

func TestMockMigrations(t *testing.T) {
	// Create mock migrations
	migrationsDir, cleanup, err := MockMigrations()
	if err != nil {
		t.Fatalf("Error creating mock migrations: %v", err)
	}
	defer cleanup()

	// Verify directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		t.Error("Migrations directory does not exist")
	}

	// Verify migration files exist
	v1Path := filepath.Join(migrationsDir, "V1__initial_schema.sql")
	if _, err := os.Stat(v1Path); os.IsNotExist(err) {
		t.Error("V1 migration file does not exist")
	}

	v2Path := filepath.Join(migrationsDir, "V2__add_test_user.sql")
	if _, err := os.Stat(v2Path); os.IsNotExist(err) {
		t.Error("V2 migration file does not exist")
	}

	// Verify file contents
	v1Content, err := os.ReadFile(v1Path)
	if err != nil {
		t.Fatalf("Failed to read V1 migration file: %v", err)
	}
	if string(v1Content) != initialSchemaMigration {
		t.Error("V1 migration file content does not match expected value")
	}

	v2Content, err := os.ReadFile(v2Path)
	if err != nil {
		t.Fatalf("Failed to read V2 migration file: %v", err)
	}
	if string(v2Content) != testUserMigration {
		t.Error("V2 migration file content does not match expected value")
	}

	// Test cleanup function
	cleanup()
	if _, err := os.Stat(migrationsDir); !os.IsNotExist(err) {
		t.Error("Cleanup function did not remove the directory")
	}
}
