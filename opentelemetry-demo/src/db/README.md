# Database Module for OpenTelemetry Demo

This module provides database functionality for the OpenTelemetry Demo microservices.

## Directory Structure

```
/src/db/
├── docker-compose.yml          # Docker Compose configuration for PostgreSQL
├── go.mod                      # Go module definition
├── migrations/                 # Database migrations directory
│   ├── README.md               # Migration documentation
│   └── versions/               # Versioned migration files
│       ├── V1__initial_schema.sql  # Initial schema migration
│       └── V2__add_test_user.sql   # Test user migration
├── postgres/                   # PostgreSQL client libraries
│   ├── connection.go           # Database connection utilities
│   ├── migration.go            # Migration runner implementation
│   └── users.go                # User repository implementation
├── scripts/                    # Utility scripts
│   ├── cleanup-usermanagement-tests.sh  # Script to clean up test resources
│   ├── setup-usermanagement-db.sh       # User management DB setup script
│   └── test-usermanagement-db-connection.go # Test script for DB connection
├── tools/                      # Command-line tools
│   └── migrate/                # Database migration tool
│       └── main.go             # Migration CLI implementation
├── README.md                   # This README file
└── setup.sh                    # Setup script to initialize database
```

## Features

- PostgreSQL database setup
- Connection management 
- User repository implementation
- Database schema migrations with version tracking
- Migration CLI tool

## Usage

### Starting the Database

```bash
# From the src/db directory
./setup.sh
```

This will:
1. Start a PostgreSQL container
2. Initialize the database with required tables using migrations
3. Set up a test user account

### Running Migrations Manually

```bash
# From the src/db directory
export DB_CONN="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
go run tools/migrate/main.go
```

You can also specify a different migrations directory:
```bash
go run tools/migrate/main.go -dir /path/to/migrations
```

### Connecting from Services

```go
import (
    "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/db/postgres"
)

// Create a connection using environment variable
conn, err := postgres.GetConnectionFromEnv("DB_CONN")
if err != nil {
    log.Fatalf("Failed to connect to database: %v", err)
}
defer conn.Close()

// Use the user repository
userRepo := postgres.NewUserRepository(conn)

// Check if a username exists
exists, err := userRepo.CheckUsername("testuser")

// Get user by username
user, err := userRepo.GetUserByUsername("testuser")

// Create a new user
userID, err := userRepo.CreateUser("newuser", "hashed_password")
```

### Testing Database Connections

```bash
# From the src/db directory
export DB_CONN="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
go run scripts/test-usermanagement-db-connection.go
```

This script will:
1. Connect to the database
2. Create a test user
3. Query for the test user
4. Clean up by removing the test user when finished

### Cleaning Up Test Resources

After running tests or when you're finished with development, you can clean up any test resources:

```bash
# From the src/db directory
./scripts/cleanup-usermanagement-tests.sh
```

This script will:
1. Delete all test users created during testing (any username starting with 'testuser_')
2. Optionally stop the PostgreSQL container when running in interactive mode
3. Show a count of remaining users in the database

## Database Migrations

Migrations follow a versioned approach:

1. Each migration is in a separate file with a version number: `V1__description.sql`
2. Migrations contain both UP (apply) and DOWN (rollback) sections
3. Migrations are tracked in a `schema_migrations` table in the database
4. Only new migrations are applied to existing databases

For more details, see the [migrations README](migrations/README.md).

## Environment Variables

The database module uses the following environment variables:

- `POSTGRES_USER`: PostgreSQL username (default: postgres)
- `POSTGRES_PASSWORD`: PostgreSQL password (default: postgres)
- `POSTGRES_DB`: PostgreSQL database name (default: postgres)
- `POSTGRES_PORT`: PostgreSQL port (default: 5432)
- `POSTGRES_HOST`: PostgreSQL host address (default: postgres)
- `DB_CONN`: Database connection string

### Connection String Format

The connection string follows this format:
```
postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable
```

### Integration with Services

For development:
```bash
# For local development
export DB_CONN="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
```

For production in the OpenTelemetry Demo:
```
# In the .env file
USER_MANAGEMENT_DB_CONN=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable
```

Services like the User Management Service expect the connection string to be passed as the `DB_CONN` environment variable. In the demo's docker-compose setup, the `.env` variable `USER_MANAGEMENT_DB_CONN` is mapped to the service's `DB_CONN` environment variable.

### Enabling Migrations

Migrations can be enabled by setting:
```
ENABLE_MIGRATIONS=true
```

When enabled, the User Management Service will automatically run necessary database migrations on startup. 