-- Migration: 001_initial_schema
-- Description: Initial schema for User Management Service

-- Up Migration
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

-- Down Migration
-- DROP INDEX IF EXISTS idx_users_username;
-- DROP TABLE IF EXISTS users; 