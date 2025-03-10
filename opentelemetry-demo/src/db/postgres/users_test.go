package postgres

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Mock MigrationRunner to avoid real migrations in tests
type mockMigrationRunner struct {
	shouldFail        bool
	invalidPathErr    bool
	invalidFileErr    bool
	migrationErr      bool
	migrationsApplied bool
}

func (m *mockMigrationRunner) ApplyMigrations() error {
	if m.shouldFail {
		if m.invalidPathErr {
			return errors.New("migrations directory not found")
		}
		if m.invalidFileErr {
			return errors.New("invalid migration file")
		}
		if m.migrationErr {
			return errors.New("database error during migration")
		}
	}
	m.migrationsApplied = true
	return nil
}

// Store original NewMigrationRunner function to restore after tests
var originalNewMigrationRunner func(db *sql.DB, migrationsDir string) MigrationRunnerInterface

func init() {
	// Save the original function so we can restore it
	originalNewMigrationRunner = NewMigrationRunner
}

func TestNewUserRepository(t *testing.T) {
	t.Run("ValidConnection", func(t *testing.T) {
		// Test repository creation with valid connection
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		conn := &Connection{DB: db}
		repo := NewUserRepository(conn)

		if repo == nil {
			t.Error("Expected repository, got nil")
		}
		if repo.conn != conn {
			t.Error("Repository connection does not match expected connection")
		}
	})

	t.Run("NilConnection", func(t *testing.T) {
		// Test with nil connection
		repo := NewUserRepository(nil)
		if repo == nil {
			t.Error("Expected repository even with nil connection, got nil")
		}
		if repo.conn != nil {
			t.Error("Expected nil connection in repository")
		}
	})
}

func TestCheckUsername(t *testing.T) {
	t.Run("ExistingUsername", func(t *testing.T) {
		// Create mock DB
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up test case
		username := "existinguser"

		// Set up expectations
		rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(username).
			WillReturnRows(rows)

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test with existing username
		exists, err := repo.CheckUsername(username)

		// Assertions
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NonExistentUsername", func(t *testing.T) {
		// Create mock DB
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up test case
		username := "nonexistentuser"

		// Set up expectations
		rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(username).
			WillReturnRows(rows)

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test with non-existent username
		exists, err := repo.CheckUsername(username)

		// Assertions
		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Create mock DB
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up test case
		username := "testuser"
		dbErr := errors.New("database error")

		// Set up expectations
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(username).
			WillReturnError(dbErr)

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test with database error during query
		exists, err := repo.CheckUsername(username)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("SuccessfulCreation", func(t *testing.T) {
		// Create mock DB
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up test case
		username := "newuser"
		passwordHash := "hashedpassword123"
		expectedID := int64(1)

		// Set up expectations
		rows := sqlmock.NewRows([]string{"id"}).AddRow(expectedID)
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(username, passwordHash).
			WillReturnRows(rows)

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test successful user creation
		id, err := repo.CreateUser(username, passwordHash)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DuplicateUsername", func(t *testing.T) {
		// Create mock DB
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up test case
		username := "existinguser"
		passwordHash := "hashedpassword123"

		// Set up expectations with a unique violation error
		duplicateErr := errors.New("pq: duplicate key value violates unique constraint")
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(username, passwordHash).
			WillReturnError(duplicateErr)

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test with duplicate username
		id, err := repo.CreateUser(username, passwordHash)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("InvalidData", func(t *testing.T) {
		// Create mock DB
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up test case - empty username should cause error
		username := ""
		passwordHash := "hashedpassword123"

		// Set up expectations with a validation error
		invalidErr := errors.New("pq: value too long for type character varying")
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(username, passwordHash).
			WillReturnError(invalidErr)

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test with invalid data
		id, err := repo.CreateUser(username, passwordHash)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Create mock DB
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up test case
		username := "newuser"
		passwordHash := "hashedpassword123"
		dbErr := errors.New("database connection error")

		// Set up expectations with a database error
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(username, passwordHash).
			WillReturnError(dbErr)

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test with database error
		id, err := repo.CreateUser(username, passwordHash)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
		assert.Equal(t, int64(0), id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUserByID(t *testing.T) {
	t.Run("ExistingUser", func(t *testing.T) {
		// Create mock DB
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up test case
		userID := int64(1)
		now := time.Now()

		// Set up expectations
		rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "created_at", "updated_at"}).
			AddRow(userID, "testuser", "hashedpw", now, now)
		mock.ExpectQuery("SELECT (.+) FROM users WHERE id = \\$1").
			WithArgs(userID).
			WillReturnRows(rows)

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test retrieving existing user
		user, err := repo.GetUserByID(userID)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userID, user.ID)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "hashedpw", user.PasswordHash)
		assert.Equal(t, now, user.CreatedAt)
		assert.Equal(t, now, user.UpdatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NonExistentUser", func(t *testing.T) {
		// Create mock DB
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up test case
		userID := int64(999)

		// Set up expectations with no rows found
		mock.ExpectQuery("SELECT (.+) FROM users WHERE id = \\$1").
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test retrieving non-existent user
		user, err := repo.GetUserByID(userID)

		// Assertions
		assert.NoError(t, err) // No error for non-existent user
		assert.Nil(t, user)    // User should be nil
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("InvalidID", func(t *testing.T) {
		// Create mock DB
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up test case - negative ID
		userID := int64(-1)

		// For simplicity, we'll return a database error for invalid ID
		dbErr := errors.New("invalid ID format")
		mock.ExpectQuery("SELECT (.+) FROM users WHERE id = \\$1").
			WithArgs(userID).
			WillReturnError(dbErr)

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test with invalid ID
		user, err := repo.GetUserByID(userID)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Create mock DB
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up test case
		userID := int64(1)
		dbErr := errors.New("database connection error")

		// Set up expectations with a database error
		mock.ExpectQuery("SELECT (.+) FROM users WHERE id = \\$1").
			WithArgs(userID).
			WillReturnError(dbErr)

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test with database error
		user, err := repo.GetUserByID(userID)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUserByUsername(t *testing.T) {
	t.Run("ExistingUser", func(t *testing.T) {
		// Create mock DB
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up test case
		username := "testuser"
		userID := int64(1)
		now := time.Now()

		// Set up expectations
		rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "created_at", "updated_at"}).
			AddRow(userID, username, "hashedpw", now, now)
		mock.ExpectQuery("SELECT (.+) FROM users WHERE username = \\$1").
			WithArgs(username).
			WillReturnRows(rows)

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test retrieving existing user
		user, err := repo.GetUserByUsername(username)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userID, user.ID)
		assert.Equal(t, username, user.Username)
		assert.Equal(t, now, user.CreatedAt)
		assert.Equal(t, now, user.UpdatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NonExistentUser", func(t *testing.T) {
		// Create mock DB
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up test case
		username := "nonexistentuser"

		// Set up expectations with no rows found
		mock.ExpectQuery("SELECT (.+) FROM users WHERE username = \\$1").
			WithArgs(username).
			WillReturnError(sql.ErrNoRows)

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test retrieving non-existent user
		user, err := repo.GetUserByUsername(username)

		// Assertions
		assert.NoError(t, err) // No error for non-existent user
		assert.Nil(t, user)    // User should be nil
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Create mock DB
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Set up test case
		username := "testuser"
		dbErr := errors.New("database connection error")

		// Set up expectations with a database error
		mock.ExpectQuery("SELECT (.+) FROM users WHERE username = \\$1").
			WithArgs(username).
			WillReturnError(dbErr)

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test with database error
		user, err := repo.GetUserByUsername(username)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestEnsureTablesExist(t *testing.T) {
	// Save the original function
	originalFunc := NewMigrationRunner
	defer func() {
		NewMigrationRunner = originalFunc
	}()

	t.Run("SuccessfulMigration", func(t *testing.T) {
		// Create mock DB
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Create mock migration runner
		mockRunner := &mockMigrationRunner{
			shouldFail:        false,
			migrationsApplied: false,
		}

		// Override migration runner creation
		NewMigrationRunner = func(db *sql.DB, migrationsDir string) MigrationRunnerInterface {
			return mockRunner
		}

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test successful migration application
		err = repo.EnsureTablesExist()

		// Assertions
		assert.NoError(t, err)
		assert.True(t, mockRunner.migrationsApplied)
	})

	t.Run("MigrationPathNotFound", func(t *testing.T) {
		// Create mock DB
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Create mock migration runner with path error
		mockRunner := &mockMigrationRunner{
			shouldFail:        true,
			invalidPathErr:    true,
			migrationsApplied: false,
		}

		// Override migration runner creation
		NewMigrationRunner = func(db *sql.DB, migrationsDir string) MigrationRunnerInterface {
			return mockRunner
		}

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test with migration path not found
		err = repo.EnsureTablesExist()

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "directory not found")
		assert.False(t, mockRunner.migrationsApplied)
	})

	t.Run("InvalidMigrationFile", func(t *testing.T) {
		// Create mock DB
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Create mock migration runner with invalid file error
		mockRunner := &mockMigrationRunner{
			shouldFail:        true,
			invalidFileErr:    true,
			migrationsApplied: false,
		}

		// Override migration runner creation
		NewMigrationRunner = func(db *sql.DB, migrationsDir string) MigrationRunnerInterface {
			return mockRunner
		}

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test with invalid migration file
		err = repo.EnsureTablesExist()

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid migration file")
		assert.False(t, mockRunner.migrationsApplied)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		// Create mock DB
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error creating mock DB: %v", err)
		}
		defer db.Close()

		// Create mock migration runner with database error
		mockRunner := &mockMigrationRunner{
			shouldFail:        true,
			migrationErr:      true,
			migrationsApplied: false,
		}

		// Override migration runner creation
		NewMigrationRunner = func(db *sql.DB, migrationsDir string) MigrationRunnerInterface {
			return mockRunner
		}

		// Create repository
		repo := &UserRepository{conn: &Connection{DB: db}}

		// Test with database error
		err = repo.EnsureTablesExist()

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.False(t, mockRunner.migrationsApplied)
	})
}
