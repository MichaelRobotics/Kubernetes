-- Migration: V3__add_email_to_user.sql
-- Description: Adds an email column to the users table
-- Services: User Management Service

-- ==================== UP MIGRATION ====================

-- Add email column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR(255) UNIQUE;

-- Add comment for the new column
COMMENT ON COLUMN users.email IS 'User email address for notifications and account recovery';

-- ==================== DOWN MIGRATION ====================

-- To roll back this migration, uncomment and execute this statement:
-- ALTER TABLE users DROP COLUMN IF EXISTS email; 