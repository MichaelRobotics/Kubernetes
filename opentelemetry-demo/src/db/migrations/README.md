# Database Migrations

This directory contains database migrations for the OpenTelemetry Demo.

## Migration Structure

Migrations are organized as follows:

- `versions/` - Contains versioned migration files
  - `V1__initial_schema.sql` - Initial database schema
  - `V2__add_test_user.sql` - Adds test user data
  - (additional migrations will be added here)

## Naming Convention

Migrations follow a versioned naming convention:

- `V{number}__{description}.sql` (e.g., `V1__initial_schema.sql`) 
- Numbers are sequential and indicate the order of execution
- Descriptions are lowercase with underscores

## Migration Format

Each migration file contains:

1. A header with metadata
2. UP migration section (changes to apply)
3. DOWN migration section (how to roll back)

Example:

```sql
-- Migration: V1__some_name.sql
-- Description: What this migration does
-- Services: Which services use this migration

-- ==================== UP MIGRATION ====================
-- SQL statements to apply the migration

-- ==================== DOWN MIGRATION ====================
-- SQL statements to roll back the migration (commented out)
```

## Running Migrations

Migrations are automatically applied when the PostgreSQL container starts. The Docker Compose configuration mounts this directory to `/docker-entrypoint-initdb.d` which PostgreSQL automatically executes on initialization.

For new databases, migrations are applied in alphanumeric order. For existing databases, be sure to run only the new migration files.

## Adding a New Migration

To add a new migration:

1. Create a new file in the `versions/` directory with the next sequential number
2. Follow the naming convention and format described above
3. Test the migration on a development database
4. Include both UP and DOWN migrations

## Rolling Back Migrations

To roll back a migration:

1. Uncomment the statements in the DOWN migration section
2. Execute them manually or use a migration tool
3. Be careful with destructive operations! 