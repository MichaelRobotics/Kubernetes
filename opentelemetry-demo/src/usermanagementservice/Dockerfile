# Build stage
FROM golang:1.23 as builder

# Create the workspace structure
WORKDIR /workspace
RUN mkdir -p src/usermanagementservice src/db/postgres

# Copy the source directories
COPY ./src/usermanagementservice/ ./src/usermanagementservice/
COPY ./src/db/postgres/connection.go ./src/db/postgres/
COPY ./src/db/postgres/go.mod ./src/db/postgres/go.sum ./src/db/postgres/

# Set working directory to usermanagementservice
WORKDIR /workspace/src/usermanagementservice

# Download dependencies
RUN go mod download

# Add grpc_health_probe as a dependency

RUN go get github.com/grpc-ecosystem/grpc-health-probe@latest
# Build the grpc_health_probe binary explicitly
RUN CGO_ENABLED=0 GOOS=linux go build -o /workspace/bin/grpc_health_probe github.com/grpc-ecosystem/grpc-health-probe

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -o /usermanagementservice

# Production stage
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy the main service binary
COPY --from=builder /usermanagementservice /usermanagementservice

# Copy the grpc_health_probe binary for health checks
COPY --from=builder /workspace/bin/grpc_health_probe /bin/grpc_health_probe

# Environment variables
ENV USER_SVC_URL=8082
ENV JWT_SECRET=your-secure-jwt-secret
ENV OTEL_EXPORTER_OTLP_ENDPOINT=otelcol:4317
ENV DB_CONN=postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable

EXPOSE 8082

ENTRYPOINT ["/usermanagementservice"]