#!/bin/bash

# This script sets up the PostgreSQL database for the OpenTelemetry Demo

# Set default environment variables
export POSTGRES_USER=${POSTGRES_USER:-postgres}
export POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-postgres}
export POSTGRES_DB=${POSTGRES_DB:-postgres}
export POSTGRES_PORT=${POSTGRES_PORT:-5432}

# Start PostgreSQL container if not already running
if ! docker ps | grep otel-postgres > /dev/null; then
    echo "Starting PostgreSQL container..."
    docker-compose up -d
else
    echo "PostgreSQL container is already running."
fi

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
for i in {1..30}; do
    if docker exec otel-postgres pg_isready -U $POSTGRES_USER > /dev/null 2>&1; then
        echo "PostgreSQL is ready."
        break
    fi
    echo "Waiting for PostgreSQL to start... ($i/30)"
    sleep 1
    if [ $i -eq 30 ]; then
        echo "Timed out waiting for PostgreSQL to start."
        exit 1
    fi
done

# Set environment variable for services
echo "Database setup complete."
echo "You can connect to the database with:"
echo "export DB_CONN=\"postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable\""

# Test connection if psql is available
if command -v psql > /dev/null; then
    echo "Testing connection to PostgreSQL..."
    PGPASSWORD=$POSTGRES_PASSWORD psql -h localhost -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT 'Connection successful!' as message;"
else
    echo "PostgreSQL client (psql) not found. Skipping connection test."
fi

echo "PostgreSQL database setup complete!" 