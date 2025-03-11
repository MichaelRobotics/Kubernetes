package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/db/postgres"
)

func main() {
	// Set up an environment variable for the connection string
	os.Setenv("DB_CONNECTION", "postgres://postgres:postgres@localhost:5432/demo?sslmode=disable")

	// Establish a connection using the environment variable
	conn, err := postgres.GetConnectionFromEnv("DB_CONNECTION")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	fmt.Println("Successfully connected to the database!")

	// Initialize the user repository
	userRepo := postgres.NewUserRepository(conn)

	// Ensure database tables exist
	if err := userRepo.EnsureTablesExist(); err != nil {
		log.Fatalf("Failed to ensure tables exist: %v", err)
	}
	fmt.Println("Database schema is up to date!")

	// Example of checking if a user exists
	username := "example_user"
	exists, err := userRepo.CheckUsername(username)
	if err != nil {
		log.Fatalf("Failed to check username: %v", err)
	}

	if exists {
		fmt.Printf("User '%s' already exists\n", username)

		// Get user by username
		user, err := userRepo.GetUserByUsername(username)
		if err != nil {
			log.Fatalf("Failed to get user: %v", err)
		}
		fmt.Printf("User ID: %d, Username: %s, Created at: %v\n",
			user.ID, user.Username, user.CreatedAt)
	} else {
		fmt.Printf("User '%s' does not exist, creating...\n", username)

		// Create a new user
		passwordHash := "hashed_password_example" // In a real app, you'd use bcrypt or similar
		userID, err := userRepo.CreateUser(username, passwordHash)
		if err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
		fmt.Printf("Created user with ID: %d\n", userID)

		// Get the user we just created
		user, err := userRepo.GetUserByID(userID)
		if err != nil {
			log.Fatalf("Failed to get user by ID: %v", err)
		}
		fmt.Printf("Retrieved user: ID: %d, Username: %s, Created at: %v\n",
			user.ID, user.Username, user.CreatedAt)
	}
}
