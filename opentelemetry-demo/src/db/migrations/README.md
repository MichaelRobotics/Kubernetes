# Database Migrations

PostgreSQL migrations for the OpenTelemetry Demo, automatically applied during container startup.

## Overview

This directory contains database schema migrations for the OpenTelemetry Demo's PostgreSQL
database. Migrations are SQL scripts that create or modify the database schema and are applied
automatically when the PostgreSQL container starts.

## How Migrations Work

1. **Automatic Execution**: 
   - The `./versions` directory is mounted to `/docker-entrypoint-initdb.d` in the PostgreSQL container
   - PostgreSQL automatically executes all SQL files in this directory in alphabetical order

2. **Version Tracking**:
   - Each migration has a version number: `V{number}__{description}.sql`
   - A `schema_migrations` table tracks which migrations have been applied
   - Only new migrations are applied to existing databases

## Migration Structure

Migrations are organized in the following directory structure:

```
migrations/
├── README.md               # This documentation
└── versions/               # Versioned SQL migration files
    ├── V1__initial_schema.sql  # Creates initial tables
    ├── V2__add_test_user.sql   # Adds a test user
    └── V{N}__{description}.sql # Additional migrations
```

## Schema

The current schema includes the following tables:

- **users**: Stores user accounts
  - `id`: Auto-incrementing primary key
  - `username`: Unique username
  - `password_hash`: Bcrypt-hashed passwords
  - `created_at`: Timestamp of account creation
  - `updated_at`: Timestamp of last update

- **schema_migrations**: Tracks applied migrations
  - `version`: Migration version number
  - `applied_at`: Timestamp when migration was applied

## Writing New Migrations

Each migration file follows this format:

```sql
-- Migration: V{number}__{description}.sql
-- Description: Brief explanation of what this migration does
-- Services: Which OpenTelemetry services use this schema

-- ==================== UP MIGRATION ====================
-- SQL statements to apply the changes
CREATE TABLE example (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);

-- ==================== DOWN MIGRATION ====================
-- SQL statements to revert the changes (commented out)
-- DROP TABLE example;
```

When adding a new migration:

1. Create a new file with the next sequential version number
2. Include both UP and DOWN migration sections
3. Use clear, descriptive comments
4. Test thoroughly before deploying

For more details, see examples in the `versions/` directory. 