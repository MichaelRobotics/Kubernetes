#!/bin/bash

# Function to clean up resources
cleanup() {
    echo "Cleaning up resources..."
    echo "Stopping and removing containers..."
    docker-compose down
    echo "Resources cleaned up."
}

# Set up trap to ensure cleanup happens even if script is interrupted
trap cleanup EXIT

echo "Starting usermanagementservice..."

# Check if the image exists
if [[ "$(docker images -q usermanagementservice:latest 2> /dev/null)" == "" ]]; then
  echo "Error: usermanagementservice:latest image not found."
  echo "Please build the image first with: docker build -t usermanagementservice:latest -f Dockerfile ../../.."
  exit 1
fi

# Stop any existing containers
docker-compose down

# Start the service
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