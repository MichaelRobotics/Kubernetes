#!/bin/bash

# This script runs database smoke tests from the smoke folder

# Set default environment variables
export POSTGRES_USER=${POSTGRES_USER:-postgres}
export POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-postgres}
export POSTGRES_DB=${POSTGRES_DB:-postgres}
export POSTGRES_PORT=${POSTGRES_PORT:-5432}
export DB_CONN="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"

# Run code from smoke folder
echo "Running code from smoke folder..."

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

echo "All smoke tests completed!" 