FROM golang:1.23 as builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker caching
COPY ./src/usermanagementservice/go.mod ./
COPY ./src/usermanagementservice/go.sum ./
RUN go mod download

# Install grpc_health_probe
RUN go install github.com/grpc-ecosystem/grpc-health-probe@latest

# Copy the service source code including pre-generated protobuf files
COPY ./src/usermanagementservice/ ./

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -o /usermanagementservice

# Create lightweight production image
FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=builder /usermanagementservice /usermanagementservice
# Copy grpc_health_probe for healthcheck
COPY --from=builder /go/bin/grpc_health_probe /bin/grpc_health_probe

ENV USER_SVC_URL=8082
ENV JWT_SECRET=your-secure-jwt-secret
ENV OTEL_EXPORTER_OTLP_ENDPOINT=otelcol:4317
ENV DB_CONN=postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable

EXPOSE 8082

# Add healthcheck using grpc_health_probe
HEALTHCHECK --interval=5s --timeout=3s --start-period=10s --retries=3 \
  CMD ["/bin/grpc_health_probe", "-addr=localhost:8082"]

ENTRYPOINT ["/usermanagementservice"] 