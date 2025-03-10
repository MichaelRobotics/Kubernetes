#!/bin/bash

# This script sets up the PostgreSQL database for the User Management Service

# Change to the database directory
cd "$(dirname "$0")/.."

# Set default environment variables
export POSTGRES_USER_MANAGEMENT=${POSTGRES_USER_MANAGEMENT:-usermanagement}
export POSTGRES_PASSWORD_MANAGEMENT=${POSTGRES_PASSWORD_MANAGEMENT:-usermanagement}
export POSTGRES_DB_MANAGEMENT=${POSTGRES_DB_MANAGEMENT:-usermanagement}
export POSTGRES_PORT_MANAGEMENT=${POSTGRES_PORT_MANAGEMENT:-5433}

# Start PostgreSQL container if not already running
if ! docker ps | grep usermanagement-postgres > /dev/null; then
    echo "Starting PostgreSQL container..."
    # Use docker-compose in the src/db directory
    docker-compose -f docker-compose.yml up -d usermanagement-postgres
else
    echo "PostgreSQL container is already running."
fi

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
for i in {1..30}; do
    if docker exec usermanagement-postgres pg_isready -U $POSTGRES_USER_MANAGEMENT > /dev/null 2>&1; then
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

# Set environment variable for the service
echo "Database setup complete."
echo "You can connect to the database with:"
echo "export DB_CONN=\"postgres://${POSTGRES_USER_MANAGEMENT}:${POSTGRES_PASSWORD_MANAGEMENT}@localhost:${POSTGRES_PORT_MANAGEMENT}/${POSTGRES_DB_MANAGEMENT}?sslmode=disable\""

# Test connection if psql is available
if command -v psql > /dev/null; then
    echo "Testing connection to PostgreSQL..."
    PGPASSWORD=$POSTGRES_PASSWORD_MANAGEMENT psql -h localhost -p $POSTGRES_PORT_MANAGEMENT -U $POSTGRES_USER_MANAGEMENT -d $POSTGRES_DB_MANAGEMENT -c "SELECT 'Connection successful!' as message;"
else
    echo "PostgreSQL client (psql) not found. Skipping connection test."
fi

echo "PostgreSQL database setup complete!"
echo "You can now run the User Management Service with the database." 