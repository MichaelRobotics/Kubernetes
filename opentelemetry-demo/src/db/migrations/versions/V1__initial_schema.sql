-- Migration: V1__initial_schema.sql
-- Description: Initial database schema for the OpenTelemetry Demo
-- Services: User Management Service

-- ==================== UP MIGRATION ====================

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indices
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Add comments
COMMENT ON TABLE users IS 'Stores user credentials and account information';
COMMENT ON COLUMN users.id IS 'Auto-incrementing primary key';
COMMENT ON COLUMN users.username IS 'Unique username for login identification';
COMMENT ON COLUMN users.password_hash IS 'Bcrypt hash of the user password';
COMMENT ON COLUMN users.created_at IS 'When the user account was created';
COMMENT ON COLUMN users.updated_at IS 'When the user account was last updated';

-- ==================== DOWN MIGRATION ====================

-- To roll back this migration, uncomment and execute these statements:
-- DROP INDEX IF EXISTS idx_users_username;
-- DROP TABLE IF EXISTS users; 