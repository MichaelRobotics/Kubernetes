#!/bin/bash

# Function to clean up resources
cleanup() {
    echo "Cleaning up resources..."
    echo "Stopping and removing containers..."
    docker-compose down
    echo "Stopping and removing database containers..."
    (cd ../db && docker-compose down)
    echo "Resources cleaned up."
}

# Set up trap to ensure cleanup happens even if script is interrupted
trap cleanup EXIT

echo "Starting usermanagementservice..."


set -a
source ../../.env
set +a

set -a
source ../../.env.override
set +a

# Stop any existing containers
docker-compose down 

# Create and start the database first
echo "Setting up database..."
(cd ../db && docker-compose up -d)

# Wait for database to be ready
echo "Waiting for database to be healthy..."
for i in {1..15}; do
  health_status=$(docker inspect --format='{{.State.Health.Status}}' postgres 2>/dev/null)
  
  if [ "$health_status" = "healthy" ]; then
    echo "Database is healthy!"
    break
  fi
  
  echo "Waiting for database to be healthy... (attempt $i/15)"
  sleep 2
done

# Start the usermanagement service
docker-compose up -d

# Check if the service is healthy
echo "Waiting for service to be healthy..."
for i in {1..30}; do
  health_status=$(docker inspect --format='{{.State.Health.Status}}' usermanagementservice 2>/dev/null)
  
  if [ "$health_status" = "healthy" ]; then
    echo "usermanagementservice is healthy!"
    echo "The service is running on http://localhost:8082"
    echo "Tests completed."
    echo "Shutting down containers..."
    # Note: We don't need to explicitly call docker-compose down here
    # as the cleanup function will be triggered by the EXIT trap
    exit 0
  fi
  
  echo "Waiting for service to be healthy... (attempt $i/30)"
  sleep 2
done

echo "Service health check timed out. Check the logs for errors:"
docker logs usermanagementservice

# No need to explicitly call cleanup as it will be triggered by the trap
# when the script exits 