# Build stage
FROM golang:1.23-bullseye AS builder

# Create the workspace structure
WORKDIR /workspace

# Create the required module directory structure
RUN mkdir -p /workspace/src/usermanagementservice /workspace/src/db/postgres /workspace/bin /workspace/src/usermanagementservice/genproto/oteldemo

# Copy go.mod and go.sum for all modules
COPY ./src/usermanagementservice/go.mod ./src/usermanagementservice/go.sum /workspace/src/usermanagementservice/
COPY ./src/db/postgres/go.mod ./src/db/postgres/go.sum /workspace/src/db/postgres/
COPY ./src/usermanagementservice/genproto/oteldemo/go.mod ./src/usermanagementservice/genproto/oteldemo/go.sum /workspace/src/usermanagementservice/genproto/oteldemo/

# Copy the source code
COPY ./src/usermanagementservice/ /workspace/src/usermanagementservice/
COPY ./src/db/postgres/ /workspace/src/db/postgres/

# Set working directory to genproto/oteldemo module first to download its dependencies
WORKDIR /workspace/src/usermanagementservice/genproto/oteldemo
RUN go mod download

# Set working directory to db/postgres module and download dependencies
WORKDIR /workspace/src/db/postgres
RUN go mod download

# Set working directory to usermanagementservice and download dependencies
WORKDIR /workspace/src/usermanagementservice
RUN go mod download

# Create a go.work file to handle all local module dependencies
RUN echo "go 1.22.1" > /workspace/go.work && \
    echo "use (" >> /workspace/go.work && \
    echo "    ./src/usermanagementservice" >> /workspace/go.work && \
    echo "    ./src/db/postgres" >> /workspace/go.work && \
    echo "    ./src/usermanagementservice/genproto/oteldemo" >> /workspace/go.work && \
    echo ")" >> /workspace/go.work

# Build the health probe with a pinned version
RUN go install github.com/grpc-ecosystem/grpc-health-probe@v0.4.19 && \
    cp $(go env GOPATH)/bin/grpc-health-probe /workspace/bin/grpc_health_probe

# Build the main service
WORKDIR /workspace
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /usermanagementservice ./src/usermanagementservice

# Production stage
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy the main service binary
COPY --from=builder /usermanagementservice /usermanagementservice

# Copy the grpc_health_probe binary for health checks
COPY --from=builder /workspace/bin/grpc_health_probe /bin/grpc_health_probe

# Environment variables
ENV USER_SVC_URL=8082
# JWT_SECRET is now expected to be passed at runtime or through Docker secrets
ENV OTEL_EXPORTER_OTLP_ENDPOINT=otelcol:4317
ENV DB_CONN=postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable

EXPOSE 8082

ENTRYPOINT ["/usermanagementservice"]