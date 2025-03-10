package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

	// Verify the connection works
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	// Verify the users table exists and has the correct schema
	var tableExists bool
	err = db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&tableExists)
	if err != nil {
		log.Fatalf("Failed to check if users table exists: %v", err)
	}

	if !tableExists {
		log.Fatal("Users table does not exist")
	}

	fmt.Println("Users table exists with the following schema:")

	// Get table schema
	rows, err := db.Query(`
		SELECT column_name, data_type, character_maximum_length
		FROM information_schema.columns
		WHERE table_name = 'users'
		ORDER BY ordinal_position
	`)
	if err != nil {
		log.Fatalf("Failed to get schema: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var columnName, dataType string
		var maxLength sql.NullInt64
		if err := rows.Scan(&columnName, &dataType, &maxLength); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		if maxLength.Valid {
			fmt.Printf("  - %s: %s(%d)\n", columnName, dataType, maxLength.Int64)
		} else {
			fmt.Printf("  - %s: %s\n", columnName, dataType)
		}
	}

	// Check how many users are in the database
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to count users: %v", err)
	}
	fmt.Printf("Total users in database: %d\n", count)

	fmt.Println("Database connection test completed successfully!")
}
