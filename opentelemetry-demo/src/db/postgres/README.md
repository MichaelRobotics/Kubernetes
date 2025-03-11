# PostgreSQL Package

Go library for PostgreSQL database interactions in the OpenTelemetry Demo application.

## Overview

This package provides utilities for connecting to PostgreSQL databases and performing
user management operations. It's designed for use by the User Management Service
and other services that require database access.

## Features

- **Connection Management**
  - Create connections from connection strings or environment variables
  - Connection pooling and lifecycle management
  - Error handling and reconnection logic
  
- **User Repository**
  - Create, retrieve, and verify user accounts
  - Password hash verification
  - Username availability checking
  
- **Database Migrations**
  - Schema creation and updates
  - Migration version tracking
  - Automatic migration application

## Usage

### Establishing a Connection

```go
// Create a connection from a connection string
conn, err := postgres.NewConnection("postgres://username:password@postgres:5432/dbname?sslmode=disable")
if err != nil {
    log.Fatalf("Failed to connect: %v", err)
}
defer conn.Close()

// Alternatively, create a connection from an environment variable
conn, err := postgres.GetConnectionFromEnv("DB_CONN")
if err != nil {
    log.Fatalf("Failed to connect: %v", err)
}
defer conn.Close()
```

### User Operations

```go
// Create a user repository
userRepo := postgres.NewUserRepository(conn)

// Check if a username exists
exists, err := userRepo.CheckUsername("testuser")

// Create a new user
userId, err := userRepo.CreateUser("newuser", "hashedpassword")

// Get a user by ID
user, err := userRepo.GetUserByID(userId)

// Get a user by username
user, err := userRepo.GetUserByUsername("newuser")

// Verify a user's password (returns true if password matches)
matched, err := userRepo.VerifyPassword("username", "password")
```

### Migrations

```go
// Ensure tables exist (applies migrations if necessary)
err := userRepo.EnsureTablesExist()
```

For more information on database migrations, see the [migrations README](../migrations/README.md).
