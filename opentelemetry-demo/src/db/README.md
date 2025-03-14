# PostgreSQL Database for OpenTelemetry Demo

PostgreSQL database service with automatic migrations for the OpenTelemetry
Demo's user management.

## Quick Start

```bash
# Start database (creates opentelemetry-demo network)
docker-compose up -d
```

## Network Configuration

The database service creates the `opentelemetry-demo` network which is shared
with the `usermanagementservice`. This network allows containers to communicate
using their service names as hostnames.

## Connection String

Connect to the database using:

```sql
postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable
```

For local development:

```bash
export DB_CONN="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
```

The hostname `postgres` resolves via Docker networking when services are on the
same network.

## Environment Variables

- `POSTGRES_USER`: PostgreSQL username (default: postgres)
- `POSTGRES_PASSWORD`: PostgreSQL password (default: postgres)
- `POSTGRES_DB`: PostgreSQL database name (default: postgres)
- `POSTGRES_PORT`: PostgreSQL port (default: 5432)
- `POSTGRES_HOST`: PostgreSQL host address (default: postgres)
- `DB_CONN`: Full database connection string

## Key Features

- **Containerized PostgreSQL**: Ready-to-use database server
- **Automatic Migrations**: Database schema applied at startup
- **User Management**: Authentication and account storage
- **Docker Network Integration**: Seamless service communication

## Directory Structure

```text
/src/db/
├── docker-compose.yml          # PostgreSQL configuration
├── migrations/                 # Database schema files
│   └── versions/               # Versioned migrations
├── postgres/                   # Client libraries
│   ├── connection.go           # Connection utilities
│   ├── migration.go            # Migration runner
│   └── users.go                # User repository
├── scripts/                    # Utility scripts
└── tools/                      # Command-line tools
```

## Testing

Run the database test script to verify connectivity:

```bash
./db-setup-test.sh
```

This will:

1. Start the database container
2. Verify database connection
3. Run smoke tests
4. Clean up automatically upon completion

For more details, see the [migrations README](migrations/README.md).
