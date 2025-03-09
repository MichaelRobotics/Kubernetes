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
├── genproto/         # Generated protobuf code
│   └── oteldemo/     # Generated service definitions
├── tests/            # Test files
│   └── mocks/        # Mock implementations for testing
│       └── db.go     # Database mock
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

1. Connect to PostgreSQL for persistent storage
2. Send telemetry data to the OpenTelemetry Collector
3. Provide authentication services to other components via gRPC

## Environment Variables

The following environment variables are required:

- `OTEL_EXPORTER_OTLP_ENDPOINT`: OpenTelemetry Collector endpoint (default: "otel-collector:4317")
- `DB_CONN`: PostgreSQL connection string
- `JWT_SECRET`: Secret key for signing JWT tokens
- `USER_SVC_URL`: Service URL and port (default: ":8080")

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
message HealthRequest {}
```

**Response:**
```protobuf
message HealthResponse {
    string status = 1;
}
```

## Generating protobuf code

The protobuf code can be generated using the following command:

```bash
cd opentelemetry-demo/src/usermanagementservice
go generate ./...
```

This will generate the necessary gRPC code in the `genproto/oteldemo` directory.
