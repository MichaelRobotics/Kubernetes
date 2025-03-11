#!/bin/bash
set -e

# Build the Docker image
cd ../..
echo "Building usermanagementservice Docker image..."
docker build -t usermanagementservice -f src/usermanagementservice/Dockerfile .

echo "Docker image built successfully!"
echo "To run the image use: docker run -p 8082:8082 usermanagementservice" 