-- Migration: V2__add_test_user.sql
-- Description: Adds a test user for development and testing purposes
-- Services: User Management Service

-- ==================== UP MIGRATION ====================

-- Create test user (password: testpassword123)
-- Only insert if the user doesn't already exist
INSERT INTO users (username, password_hash) 
VALUES ('testuser', '$2a$10$QGfO0JUVuG5R.lQGXSIzd.pBB7WmJjkJ6zf6jE/oyGqhR8tGWRYMG')
ON CONFLICT (username) DO NOTHING;

-- ==================== DOWN MIGRATION ====================

-- To roll back this migration, uncomment and execute this statement:
-- DELETE FROM users WHERE username = 'testuser'; 