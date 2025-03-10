# Centralized Database Module for OpenTelemetry Demo

This module provides centralized database management for the OpenTelemetry Demo microservices.

## Directory Structure

```
/src/db/
├── docker-compose.yml          # Docker Compose configuration for all databases
├── go.mod                      # Go module definition
├── init/                       # Database initialization scripts
│   └── usermanagement-init.sql # Init script for User Management Service
├── migrations/                 # Database migration scripts
│   └── usermanagement/         # Migrations for User Management Service
│       └── 001_initial_schema.sql  
├── postgres/                   # PostgreSQL client libraries
│   ├── connection.go           # Database connection utilities
│   └── users.go                # User repository implementation
├── README.md                   # This README file
└── scripts/                    # Database management scripts
    ├── setup-usermanagement-db.sh             # Script to set up User Management database
    └── test-usermanagement-db-connection.go   # Test connection to User Management database
```

## Features

- Centralized database configuration
- PostgreSQL connection management
- Environment variable-based configuration
- Service-specific repositories
- Docker Compose setup for development

## Usage

### Starting the Database Containers

```bash
# From the src/db directory
docker-compose up -d
```

### Setting Up a Specific Service Database

```bash
# For User Management Service
./scripts/setup-usermanagement-db.sh
```

### Testing Database Connections

```bash
# For User Management Service
export DB_CONN="postgres://usermanagement:usermanagement@localhost:5433/usermanagement?sslmode=disable"
go run scripts/test-usermanagement-db-connection.go
```

## Environment Variables

The following environment variables can be set in the `.env` file at the root of the project:

### Global Database Configuration
```
POSTGRES_IMAGE=postgres:17.3    # PostgreSQL Docker image to use
```

### User Management Service Database
```
# User Management Service Database
POSTGRES_USER_MANAGEMENT=usermanagement
POSTGRES_PASSWORD_MANAGEMENT=usermanagement
POSTGRES_DB_MANAGEMENT=usermanagement
POSTGRES_PORT_MANAGEMENT=5433
USER_MANAGEMENT_DB_CONN=postgres://${POSTGRES_USER_MANAGEMENT}:${POSTGRES_PASSWORD_MANAGEMENT}@${POSTGRES_HOST}:${POSTGRES_PORT_MANAGEMENT}/${POSTGRES_DB_MANAGEMENT}?sslmode=disable
```

## Adding a New Service Database

1. Create a directory for migrations: `mkdir -p migrations/newservice`
2. Add initialization SQL: `touch init/newservice-init.sql`
3. Add setup script: `cp scripts/setup-usermanagement-db.sh scripts/setup-newservice-db.sh`
4. Modify the script for your new service
5. Add service-specific environment variables to the root `.env` file
6. Add the new service to `docker-compose.yml`

## Connection from Services

Services should use the environment variable with their connection string:

```go
import (
    "github.com/open-telemetry/opentelemetry-demo/src/db/postgres"
)

// Create a connection using environment variable
conn, err := postgres.GetConnectionFromEnv("DB_CONN")
if err != nil {
    log.Fatalf("Failed to connect to database: %v", err)
}
defer conn.Close()

// Use the connection for database operations
userRepo := postgres.NewUserRepository(conn)
``` 