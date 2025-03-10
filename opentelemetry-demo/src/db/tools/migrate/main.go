package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/db/postgres"
)

func main() {
	// Define flags
	migrationsDir := flag.String("dir", filepath.Join("..", "..", "migrations", "versions"), "Directory containing migration files")
	dbConnStr := flag.String("conn", os.Getenv("DB_CONN"), "Database connection string (can also use DB_CONN env var)")
	help := flag.Bool("help", false, "Show help")

	// Parse flags
	flag.Parse()

	// Show help if requested
	if *help {
		printUsage()
		os.Exit(0)
	}

	// Check if we have a database connection string
	if *dbConnStr == "" {
		fmt.Println("Error: Database connection string is required.")
		fmt.Println("Set it using the -conn flag or the DB_CONN environment variable.")
		printUsage()
		os.Exit(1)
	}

	// Connect to the database
	fmt.Println("Connecting to database...")
	db, err := postgres.GetConnection(*dbConnStr)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create a migration runner
	fmt.Printf("Using migrations directory: %s\n", *migrationsDir)
	runner := postgres.NewMigrationRunner(db, *migrationsDir)

	// Apply migrations
	fmt.Println("Applying migrations...")
	err = runner.ApplyMigrations()
	if err != nil {
		fmt.Printf("Error applying migrations: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Migrations completed successfully!")
}

func printUsage() {
	fmt.Println("Database Migration Tool for OpenTelemetry Demo")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  migrate [options]")
	fmt.Println("")
	fmt.Println("Options:")
	flag.PrintDefaults()
}
