package postgres

import (
	"database/sql"
	"errors"
	"log"
	"path/filepath"
	"runtime"
	"time"
)

// User represents a user in the database
type User struct {
	ID           int64
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserRepository provides database operations for users
type UserRepository struct {
	conn *Connection
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(conn *Connection) *UserRepository {
	return &UserRepository{conn: conn}
}

// CheckUsername checks if a username already exists
func (r *UserRepository) CheckUsername(username string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"
	err := r.conn.DB.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(username, passwordHash string) (int64, error) {
	var id int64
	query := `
		INSERT INTO users (username, password_hash, created_at, updated_at) 
		VALUES ($1, $2, NOW(), NOW()) 
		RETURNING id
	`
	err := r.conn.DB.QueryRow(query, username, passwordHash).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id int64) (*User, error) {
	var user User
	query := "SELECT id, username, password_hash, created_at, updated_at FROM users WHERE id = $1"
	err := r.conn.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (r *UserRepository) GetUserByUsername(username string) (*User, error) {
	var user User
	query := "SELECT id, username, password_hash, created_at, updated_at FROM users WHERE username = $1"
	err := r.conn.DB.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &user, nil
}

// EnsureTablesExist ensures that the required database tables exist
// This now uses the migration system instead of direct table creation
func (r *UserRepository) EnsureTablesExist() error {
	// Get the path to the migrations directory
	_, filename, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filepath.Dir(filename))
	migrationsDir := filepath.Join(basePath, "migrations", "versions")

	log.Printf("Running migrations from: %s", migrationsDir)

	// Create a migration runner
	runner := NewMigrationRunner(r.conn.DB, migrationsDir)

	// Apply migrations
	if err := runner.ApplyMigrations(); err != nil {
		log.Printf("Failed to apply migrations: %v", err)
		return err
	}

	log.Printf("Database migrations applied successfully")
	return nil
}
