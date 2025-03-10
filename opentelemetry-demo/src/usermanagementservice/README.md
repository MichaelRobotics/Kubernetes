# Go User Management Service

This service provides user management functionality for the OpenTelemetry Demo, including user registration, login, and authentication using gRPC.

## Features

- **User Registration**: Create new user accounts with username/password
- **User Authentication**: Login with credentials and receive JWT tokens
- **Health Check**: Simple endpoint to verify service health
- **OpenTelemetry Instrumentation**: Full tracing support
- **gRPC API**: Service implements gRPC interfaces

## Code Structure

The service follows a modular design pattern:

```
usermanagementservice/
├── handlers/         # API endpoint handlers
│   ├── auth.go       # Authentication logic (register/login)
│   └── health.go     # Health check handler
├── models/           # Data models
│   └── user.go       # User data structure
├── genproto/         # Generated protobuf code
│   └── oteldemo/     # Generated service definitions
├── tests/            # Test files
│   └── mocks/        # Mock implementations for testing
├── Dockerfile        # Container definition
├── go.mod            # Go module dependencies
├── go.sum            # Go module checksums
├── main.go           # Application entry point and gRPC server
├── tools.go          # Tools for protobuf generation
└── README.md         # This file
```

## Testing

Tests for this service can be run using the Makefile targets defined in the root OpenTelemetry Demo Makefile:

```bash
# Run all tests
make test-usermanagementservice

# Run unit tests only
make test-usermanagementservice-unit

# Run tests with coverage report
make test-usermanagementservice-coverage
```

The coverage report will be generated at `coverage.html` in the service directory.

## CI/CD Integration

This service includes GitHub Actions integration that automatically runs tests when changes are pushed to files in the usermanagementservice directory. The workflow file is located at `.github/workflows/test-usermanagementservice.yml`.

## Integration with OpenTelemetry Demo

The service is configured to:

1. Connect to PostgreSQL for persistent storage via the centralized database module
2. Send telemetry data to the OpenTelemetry Collector
3. Provide authentication services to other components via gRPC

## Database Integration

This service uses the centralized database module located at `/src/db` for all database operations. The database module provides:

- Connection management
- User repository with CRUD operations
- Schema creation and migration

The database connection is configured using the `DB_CONN` environment variable, which should contain a complete PostgreSQL connection string:

```
DB_CONN=postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable
```

### Database Migrations

By default, automatic migrations are disabled when the service starts. To enable automatic migrations, set:

```
ENABLE_MIGRATIONS=true
```

When migrations are enabled, the service will:
1. Attempt to create required tables if they don't exist
2. Log the migration process

When migrations are disabled (default), you'll need to apply migrations manually using the migration tool:

```bash
cd opentelemetry-demo/src/db
export DB_CONN="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
go run tools/migrate/main.go
```

## Authentication

The service uses JWT (JSON Web Tokens) for authentication. When a user logs in:

1. Credentials are verified against the database
2. A JWT token is generated using the secret key
3. The token is returned to the client for future authenticated requests

The JWT token contains the user ID as a claim and is signed with the secret specified in the `JWT_SECRET` environment variable.

## Environment Variables

The following environment variables are directly used by the service:

- `DB_CONN`: PostgreSQL connection string in the format:
  ```
  postgres://username:password@host:port/database?sslmode=disable
  ```

- `JWT_SECRET`: Secret key for signing JWT tokens
  ```
  JWT_SECRET=your-secure-jwt-secret-here
  ```

- `USER_SVC_URL`: Service URL and port (default: ":8080")
  ```
  USER_SVC_URL=:8082
  ```

- `OTEL_EXPORTER_OTLP_ENDPOINT`: OpenTelemetry Collector endpoint
  ```
  OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317
  ```

Optional environment variables:

- `ENABLE_MIGRATIONS`: Set to "true" to enable automatic database migrations on startup (default: false)
  ```
  ENABLE_MIGRATIONS=true
  ```

## gRPC API

The service exposes the following gRPC endpoints:

### Register

Register a new user account.

**Request:**
```protobuf
message RegisterRequest {
    string username = 1;
    string password = 2;
}
```

**Response:**
```protobuf
message RegisterResponse {
    int32 user_id = 1;
    string username = 2;
    string message = 3;
}
```

### Login

Authenticate and receive a JWT token.

**Request:**
```protobuf
message LoginRequest {
    string username = 1;
    string password = 2;
}
```

**Response:**
```protobuf
message LoginResponse {
    string token = 1;
    int32 user_id = 2;
}
```

### Health

Check service health.

**Request:**
```protobuf
message HealthCheckRequest {
    string service = 1;
}
```

**Response:**
```protobuf
message HealthCheckResponse {
    enum ServingStatus {
        UNKNOWN = 0;
        SERVING = 1;
        NOT_SERVING = 2;
        SERVICE_UNKNOWN = 3;
    }
    ServingStatus status = 1;
}
```

## Generating protobuf code

The protobuf code can be generated using the following command:

```bash
cd opentelemetry-demo/src/usermanagementservice
go generate ./...
```

This will generate the necessary gRPC code in the `genproto/oteldemo` directory.
