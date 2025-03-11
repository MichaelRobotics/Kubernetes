#!/bin/bash

# Function to clean up resources
cleanup() {
    echo "Cleaning up resources..."
    echo "Stopping and removing containers..."
    docker-compose down -v
    echo "Resources cleaned up."
}

# Set up trap to ensure cleanup happens even if script is interrupted
trap cleanup EXIT

# Get the directory where the script is located
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
cd "$SCRIPT_DIR"

# Source environment variables from .env file
echo "Sourcing environment variables from .env file..."
set -a
source ../../.env
set +a

set -a
source ../../.env.override
set +a

# Check if the required env vars are set
if [ -z "$POSTGRES_USER" ] || [ -z "$POSTGRES_PASSWORD" ] || [ -z "$POSTGRES_DB" ] || [ -z "$POSTGRES_PORT" ]; then
    echo "One or more required environment variables are not set"
    echo "Ensure POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, and POSTGRES_PORT are set in .env"
    exit 1
fi


# Start the database
echo "Starting PostgreSQL database..."
docker-compose up -d

# Wait for the database to be ready
echo "Waiting for database to be ready..."
for i in {1..30}; do
    if docker-compose exec postgres pg_isready -h localhost -U "$POSTGRES_USER" -d "$POSTGRES_DB" > /dev/null 2>&1; then
        echo "Database is ready!"
        break
    fi
    echo "Waiting for database to be ready... (attempt $i/30)"
    sleep 2
    if [ $i -eq 30 ]; then
        echo "Database failed to start within timeout."
        exit 1
    fi
done

# Set the DB_CONN environment variable for the smoke tests
export DB_CONN="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"
echo "Using connection string: $DB_CONN"

# Run smoke tests
echo "Running smoke tests..."

# Check if smoke directory exists
if [ -d "./smoke" ]; then
    # List available Go files in smoke directory
    echo "Found the following test files in smoke directory:"
    ls -la ./smoke/*.go 2>/dev/null || echo "No Go files found in smoke directory."
    
    # Run each Go file in the smoke directory
    for test_file in ./smoke/*.go; do
        if [ -f "$test_file" ]; then
            echo "Running $test_file..."
            go run "$test_file"
            if [ $? -eq 0 ]; then
                echo "✅ $test_file completed successfully!"
            else
                echo "❌ $test_file failed!"
                exit 1
            fi
        fi
    done
else
    echo "Smoke directory not found."
    exit 1
fi

echo "All smoke tests completed successfully!"
echo "Cleaning up..."

# Cleanup function will be called automatically due to the trap 