FROM golang:1.22 as builder

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /usermanagementservice

# Create lightweight production image
FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=builder /usermanagementservice /usermanagementservice

ENV USER_SVC_URL=:8080
ENV JWT_SECRET=changeme
ENV OTEL_EXPORTER_OTLP_ENDPOINT=otelcol:4317
ENV DB_CONN=postgresql://postgres:postgres@postgres:5432/users?sslmode=disable

EXPOSE 8080

ENTRYPOINT ["/usermanagementservice"] 