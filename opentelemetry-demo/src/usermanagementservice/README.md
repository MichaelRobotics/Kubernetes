# Go User Management Service

User authentication microservice for OpenTelemetry Demo with gRPC API for registration,
login, and JWT-based authentication.

## Features

- User registration and authentication with JWT tokens
- PostgreSQL database integration
- OpenTelemetry instrumentation
- Health check endpoint
- gRPC API

## Code Structure

```
├── handlers/         # Auth and health endpoint handlers
├── models/           # User data model
├── genproto/         # Generated gRPC code
├── tests/            # Test files and mocks
└── main.go           # Service entry point and server
```

## Network Configuration

Connects to PostgreSQL via the `opentelemetry-demo` network:
- Database service creates this network automatically
- Container name `postgres` is used as hostname for connections
- Communication occurs over the shared Docker network

## Quick Start

```bash
# Start services in proper order
cd ../db && docker-compose up -d  # Start database first
cd ../usermanagementservice && docker-compose up -d
```

## Database

- Connection string: `postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:${POSTGRES_PORT}/${POSTGRES_DB}`
- Tables: `users` (accounts), `schema_migrations` (applied migrations)
- Enable migrations: `ENABLE_MIGRATIONS=true`
- Password storage: bcrypt hashed with appropriate cost factor

## Environment Variables

- `DB_CONN`: PostgreSQL connection string
- `JWT_SECRET`: Secret for signing JWT tokens
- `USER_SVC_PORT`: Service port (default: "8080")
- `OTEL_EXPORTER_OTLP_ENDPOINT`: Collector endpoint
- `OTEL_RESOURCE_ATTRIBUTES`: OpenTelemetry resource attributes

## API Endpoints

### Register
- Request: `RegisterRequest{username, password}`
- Response: `RegisterResponse{user_id, username, message}`

### Login
- Request: `LoginRequest{username, password}`
- Response: `LoginResponse{token, user_id}`

### Health
- Request: `HealthCheckRequest{service}`
- Response: `HealthCheckResponse{status}`

## Testing

```bash
# Unit tests
make test-usermanagementservice-unit

# All tests with coverage
make test-usermanagementservice-coverage
```

## Development

Generate gRPC code:
```bash
cd opentelemetry-demo/src/usermanagementservice
go generate ./...
```

For more details on database integration, see the [db README](../db/README.md).
