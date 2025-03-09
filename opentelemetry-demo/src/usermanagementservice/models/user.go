// This file is kept for reference only. The actual data structures are now defined in the protobuf files.
// See genproto/oteldemo/demo.pb.go for the generated data structures.
package models

// User represents the user model stored in the database
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // Never expose in JSON
}

// Credentials represents the login/registration request body
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response with JWT token
type LoginResponse struct {
	Token  string `json:"token"`
	UserID int    `json:"user_id"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error string `json:"error"`
}
