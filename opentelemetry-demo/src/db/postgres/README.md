# PostgreSQL Database Package

This package provides utilities for PostgreSQL database interaction within the OpenTelemetry Demo application.

## Features

- Connection management with PostgreSQL
- User management operations (create, retrieve, check)
- Database migration system

## Usage

### Establishing a Connection

```go
import "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/db/postgres"

// Create a connection from a connection string
conn, err := postgres.NewConnection("postgres://username:password@localhost:5432/dbname?sslmode=disable")
if err != nil {
    // Handle error
}
defer conn.Close()

// Or create a connection from an environment variable
conn, err := postgres.GetConnectionFromEnv("DATABASE_URL")
if err != nil {
    // Handle error
}
defer conn.Close()
```

### User Operations

```go
// Create a user repository
userRepo := postgres.NewUserRepository(conn)

// Check if a username exists
exists, err := userRepo.CheckUsername("testuser")

// Create a new user
userId, err := userRepo.CreateUser("newuser", "hashedpassword")

// Get a user by ID
user, err := userRepo.GetUserByID(userId)

// Get a user by username
user, err := userRepo.GetUserByUsername("newuser")
```

### Migrations

```go
// Ensure tables exist (applies migrations)
err := userRepo.EnsureTablesExist()
```

## Database Migrations

Migrations are placed in the `../migrations/versions` directory with the naming pattern `V{version}__{name}.sql`.

Each migration file should contain:

```sql
-- Description: Brief description of the migration

-- ==================== UP MIGRATION ====================
-- SQL statements for applying the migration

-- ==================== DOWN MIGRATION ====================
-- SQL statements for reverting the migration
```

## Development

To add a new database operation:

1. Create a new method on an existing repository or add a new repository struct
2. Write tests for the new functionality
3. Implement the method using the SQL query builder or raw SQL
4. Add corresponding migrations if needed 