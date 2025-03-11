# Build stage
FROM golang:1.23-bullseye as builder

# Create the workspace structure
WORKDIR /workspace
RUN mkdir -p src/usermanagementservice src/db/postgres

# Copy just the go.mod and go.sum files first
COPY ./src/usermanagementservice/go.mod ./src/usermanagementservice/go.sum ./src/usermanagementservice/
COPY ./src/db/postgres/go.mod ./src/db/postgres/go.sum ./src/db/postgres/

# Download dependencies before copying the rest of the code
WORKDIR /workspace/src/usermanagementservice
RUN go mod download

# Now copy the rest of the code
COPY ./src/usermanagementservice/ ./src/usermanagementservice/
COPY ./src/db/postgres/connection.go ./src/db/postgres/

# Set working directory to usermanagementservice
WORKDIR /workspace/src/usermanagementservice

# Add grpc_health_probe and build in a single layer
RUN go get github.com/grpc-ecosystem/grpc-health-probe@latest && \
    CGO_ENABLED=0 GOOS=linux go build -o /workspace/bin/grpc_health_probe github.com/grpc-ecosystem/grpc-health-probe && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /usermanagementservice

# Production stage
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy the main service binary
COPY --from=builder /usermanagementservice /usermanagementservice

# Copy the grpc_health_probe binary for health checks
COPY --from=builder /workspace/bin/grpc_health_probe /bin/grpc_health_probe

# Environment variables
ENV USER_SVC_URL=8082
ENV JWT_SECRET=INSECURE_PLACEHOLDER_CHANGE_ME
ENV OTEL_EXPORTER_OTLP_ENDPOINT=otelcol:4317
ENV DB_CONN=postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable

EXPOSE 8082

ENTRYPOINT ["/usermanagementservice"]