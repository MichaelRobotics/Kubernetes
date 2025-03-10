// This file is kept for reference only. The actual data structures are now defined in the protobuf files.
// See genproto/oteldemo/demo.pb.go for the generated data structures.
package models

import "time"

// User represents a user in the system
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Not exposed in JSON
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Credentials represents a user login/registration request
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response with JWT token
type LoginResponse struct {
	Token  string `json:"token"`
	UserID int64  `json:"user_id"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error string `json:"error"`
}
