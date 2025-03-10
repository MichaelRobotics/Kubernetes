package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Get database connection string from environment variable
	dbConn := os.Getenv("DB_CONN")
	if dbConn == "" {
		log.Fatal("DB_CONN environment variable not set")
	}

	// Connect to the database
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("Successfully connected to the database!")

	// Create a user
	username := fmt.Sprintf("testuser_%d", time.Now().Unix())
	// Commented out as not used: password := "password123"
	passwordHash := "$2a$10$QGfO0JUVuG5R.lQGXSIzd.pBB7WmJjkJ6zf6jE/oyGqhR8tGWRYMG" // hash for "testpassword123"

	_, err = db.Exec(
		"INSERT INTO users (username, password_hash) VALUES ($1, $2)",
		username, passwordHash,
	)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("Successfully created user: %s\n", username)

	// Retrieve the user
	var (
		userID                int64
		retrievedUsername     string
		retrievedPasswordHash string
		createdAt             time.Time
		updatedAt             time.Time
	)

	err = db.QueryRow(
		"SELECT id, username, password_hash, created_at, updated_at FROM users WHERE username = $1",
		username,
	).Scan(&userID, &retrievedUsername, &retrievedPasswordHash, &createdAt, &updatedAt)

	if err != nil {
		log.Fatalf("Failed to retrieve user: %v", err)
	}

	fmt.Printf("Retrieved user: ID=%d, Username=%s, CreatedAt=%s, UpdatedAt=%s\n",
		userID, retrievedUsername, createdAt.Format(time.RFC3339), updatedAt.Format(time.RFC3339))

	// Count total users
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to count users: %v", err)
	}
	fmt.Printf("Total users in database: %d\n", count)

	// Clean up: Delete the test user created during this test
	_, err = db.Exec("DELETE FROM users WHERE username = $1", username)
	if err != nil {
		log.Fatalf("Failed to delete test user: %v", err)
	}
	fmt.Printf("Successfully deleted test user: %s\n", username)

	// Verify the deletion
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		log.Fatalf("Failed to verify user deletion: %v", err)
	}

	if count == 0 {
		fmt.Println("Verified test user was deleted successfully")
	} else {
		log.Fatalf("Failed to delete test user: user still exists in database")
	}

	// Count total users after deletion
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to count users after deletion: %v", err)
	}
	fmt.Printf("Total users in database after cleanup: %d\n", count)

	fmt.Println("Database connection test completed successfully!")
}
