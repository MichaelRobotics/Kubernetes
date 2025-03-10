package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Migration represents a database migration
type Migration struct {
	Version     int
	Name        string
	Description string
	UpSQL       string
	DownSQL     string
	FilePath    string
}

// MigrationRunner manages database migrations
type MigrationRunner struct {
	DB        *sql.DB
	Directory string
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *sql.DB, migrationsDir string) *MigrationRunner {
	return &MigrationRunner{
		DB:        db,
		Directory: migrationsDir,
	}
}

// EnsureMigrationTable creates the migrations table if it doesn't exist
func (m *MigrationRunner) EnsureMigrationTable() error {
	_, err := m.DB.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

// LoadMigrations loads all migration files from the specified directory
func (m *MigrationRunner) LoadMigrations() ([]Migration, error) {
	var migrations []Migration

	// Walk through migration files
	err := filepath.Walk(m.Directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-SQL files
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".sql") {
			return nil
		}

		// Parse the filename: V1__description.sql
		filename := info.Name()
		if !strings.HasPrefix(filename, "V") || !strings.Contains(filename, "__") {
			return nil
		}

		versionStr := strings.Split(filename, "__")[0][1:] // Remove 'V' prefix
		version := 0
		_, err = fmt.Sscanf(versionStr, "%d", &version)
		if err != nil {
			return fmt.Errorf("invalid migration version in %s: %v", filename, err)
		}

		// Read the file content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Parse the file content to extract description, up and down migrations
		contentStr := string(content)
		upSQL, downSQL := splitMigration(contentStr)

		description := extractDescription(contentStr)
		name := strings.Split(filename, "__")[1]
		name = strings.TrimSuffix(name, ".sql")

		migrations = append(migrations, Migration{
			Version:     version,
			Name:        name,
			Description: description,
			UpSQL:       upSQL,
			DownSQL:     downSQL,
			FilePath:    path,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// GetAppliedMigrations returns a list of already applied migrations
func (m *MigrationRunner) GetAppliedMigrations() (map[int]bool, error) {
	applied := make(map[int]bool)

	rows, err := m.DB.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		// If the table doesn't exist, return an empty map
		if strings.Contains(err.Error(), "does not exist") {
			return applied, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, nil
}

// ApplyMigrations applies all pending migrations
func (m *MigrationRunner) ApplyMigrations() error {
	// Ensure the migration table exists
	if err := m.EnsureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %v", err)
	}

	// Load all migrations
	migrations, err := m.LoadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %v", err)
	}

	// Get applied migrations
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %v", err)
	}

	// Begin a transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Apply pending migrations
	for _, migration := range migrations {
		if applied[migration.Version] {
			fmt.Printf("Migration %d (%s) already applied, skipping\n", migration.Version, migration.Name)
			continue
		}

		fmt.Printf("Applying migration %d: %s\n", migration.Version, migration.Name)

		// Execute the migration
		_, err = tx.Exec(migration.UpSQL)
		if err != nil {
			return fmt.Errorf("failed to apply migration %d (%s): %v", migration.Version, migration.Name, err)
		}

		// Record the migration
		_, err = tx.Exec(
			"INSERT INTO schema_migrations (version, name) VALUES ($1, $2)",
			migration.Version, migration.Name,
		)
		if err != nil {
			return fmt.Errorf("failed to record migration %d: %v", migration.Version, err)
		}

		fmt.Printf("Successfully applied migration %d\n", migration.Version)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// Helper function to extract the description from a migration file
func extractDescription(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "-- Description:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "-- Description:"))
		}
	}
	return ""
}

// Helper function to split migration into up and down parts
func splitMigration(content string) (string, string) {
	parts := strings.Split(content, "-- ==================== DOWN MIGRATION ====================")

	if len(parts) < 2 {
		return content, ""
	}

	upPart := strings.Split(parts[0], "-- ==================== UP MIGRATION ====================")
	if len(upPart) < 2 {
		return parts[0], parts[1]
	}

	return upPart[1], parts[1]
}
