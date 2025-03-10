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

-- Create test user (password: testpassword123)
INSERT INTO users (username, password_hash) 
VALUES ('testuser', '$2a$10$QGfO0JUVuG5R.lQGXSIzd.pBB7WmJjkJ6zf6jE/oyGqhR8tGWRYMG')
ON CONFLICT (username) DO NOTHING;

-- Add comments
COMMENT ON TABLE users IS 'Stores user credentials and account information';
COMMENT ON COLUMN users.id IS 'Auto-incrementing primary key';
COMMENT ON COLUMN users.username IS 'Unique username for login identification';
COMMENT ON COLUMN users.password_hash IS 'Bcrypt hash of the user password';
COMMENT ON COLUMN users.created_at IS 'When the user account was created';
COMMENT ON COLUMN users.updated_at IS 'When the user account was last updated'; 