package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserModel(t *testing.T) {
	// Create a user
	user := User{
		ID:           123,
		Username:     "testuser",
		PasswordHash: "hashed_password",
	}

	// Test fields
	assert.Equal(t, 123, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "hashed_password", user.PasswordHash)
}

func TestUserJSONSerialization(t *testing.T) {
	// Create a user
	user := User{
		ID:           123,
		Username:     "testuser",
		PasswordHash: "this_should_not_be_visible",
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err)

	// Convert to string for assertions
	jsonString := string(jsonData)

	// Verify password hash is not included in JSON
	assert.Contains(t, jsonString, "123")
	assert.Contains(t, jsonString, "testuser")
	assert.NotContains(t, jsonString, "this_should_not_be_visible")
}

func TestCredentialsModel(t *testing.T) {
	// Create credentials
	creds := Credentials{
		Username: "testuser",
		Password: "password123",
	}

	// Test fields
	assert.Equal(t, "testuser", creds.Username)
	assert.Equal(t, "password123", creds.Password)
}

func TestLoginResponseModel(t *testing.T) {
	// Create login response
	response := LoginResponse{
		Token:  "jwt.token.here",
		UserID: 123,
	}

	// Test fields
	assert.Equal(t, "jwt.token.here", response.Token)
	assert.Equal(t, 123, response.UserID)
}

func TestErrorResponseModel(t *testing.T) {
	// Create error response
	errResponse := ErrorResponse{
		Error: "Authentication failed",
	}

	// Test fields
	assert.Equal(t, "Authentication failed", errResponse.Error)
}

func TestErrorResponseJSONSerialization(t *testing.T) {
	// Create error response
	errResponse := ErrorResponse{
		Error: "Authentication failed",
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(errResponse)
	assert.NoError(t, err)

	// Verify JSON structure
	expectedJSON := `{"error":"Authentication failed"}`
	assert.JSONEq(t, expectedJSON, string(jsonData))
}
