package handlers

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/opentelemetry/demo/src/usermanagementservice/genproto/oteldemo"
	"github.com/opentelemetry/demo/src/usermanagementservice/tests/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRegister_Success(t *testing.T) {
	// Setup
	db, mock, err := mocks.NewDBMock()
	assert.NoError(t, err)
	defer db.Close()

	tracer := mocks.NewMockTracer()
	jwtSecret := []byte("test-secret")
	
	handler := NewAuthHandler(db, tracer, jwtSecret)
	
	// Test data
	username := "testuser"
	password := "password123"
	userID := int32(1)
	
	// Mock expectations
	mocks.ExpectUsernameCheck(mock, username, false) // Username doesn't exist
	mocks.ExpectInsertUser(mock, username, password, userID)
	
	// Call the method
	req := &pb.RegisterRequest{
		Username: username,
		Password: password,
	}
	
	resp, err := handler.Register(context.Background(), req)
	
	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, userID, resp.UserId)
	assert.Equal(t, username, resp.Username)
	assert.Equal(t, "User registered successfully", resp.Message)
	
	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegister_UsernameTooShort(t *testing.T) {
	// Setup
	db, mock, err := mocks.NewDBMock()
	assert.NoError(t, err)
	defer db.Close()

	tracer := mocks.NewMockTracer()
	jwtSecret := []byte("test-secret")
	
	handler := NewAuthHandler(db, tracer, jwtSecret)
	
	// Test data
	username := "ab" // Too short (less than 3 characters)
	password := "password123"
	
	// Call the method
	req := &pb.RegisterRequest{
		Username: username,
		Password: password,
	}
	
	resp, err := handler.Register(context.Background(), req)
	
	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	
	// Check error code and message
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "username must be at least 3 characters")
	
	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegister_PasswordTooShort(t *testing.T) {
	// Setup
	db, mock, err := mocks.NewDBMock()
	assert.NoError(t, err)
	defer db.Close()

	tracer := mocks.NewMockTracer()
	jwtSecret := []byte("test-secret")
	
	handler := NewAuthHandler(db, tracer, jwtSecret)
	
	// Test data
	username := "testuser"
	password := "pass" // Too short (less than 8 characters)
	
	// Call the method
	req := &pb.RegisterRequest{
		Username: username,
		Password: password,
	}
	
	resp, err := handler.Register(context.Background(), req)
	
	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	
	// Check error code and message
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "password must be at least 8 characters")
	
	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegister_UsernameExists(t *testing.T) {
	// Setup
	db, mock, err := mocks.NewDBMock()
	assert.NoError(t, err)
	defer db.Close()

	tracer := mocks.NewMockTracer()
	jwtSecret := []byte("test-secret")
	
	handler := NewAuthHandler(db, tracer, jwtSecret)
	
	// Test data
	username := "testuser"
	password := "password123"
	
	// Mock expectations
	mocks.ExpectUsernameCheck(mock, username, true) // Username exists
	
	// Call the method
	req := &pb.RegisterRequest{
		Username: username,
		Password: password,
	}
	
	resp, err := handler.Register(context.Background(), req)
	
	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	
	// Check error code and message
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, st.Code())
	assert.Contains(t, st.Message(), "username already exists")
	
	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_Success(t *testing.T) {
	// Setup
	db, mock, err := mocks.NewDBMock()
	assert.NoError(t, err)
	defer db.Close()

	tracer := mocks.NewMockTracer()
	jwtSecret := []byte("test-secret")
	
	handler := NewAuthHandler(db, tracer, jwtSecret)
	
	// Test data
	username := "testuser"
	password := "password123"
	userID := int32(1)
	
	// Hash the password for storage
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	
	// Mock expectations
	mocks.ExpectGetUser(mock, username, userID, string(hashedPassword))
	
	// Call the method
	req := &pb.LoginRequest{
		Username: username,
		Password: password,
	}
	
	resp, err := handler.Login(context.Background(), req)
	
	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, userID, resp.UserId)
	
	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_UserNotFound(t *testing.T) {
	// Setup
	db, mock, err := mocks.NewDBMock()
	assert.NoError(t, err)
	defer db.Close()

	tracer := mocks.NewMockTracer()
	jwtSecret := []byte("test-secret")
	
	handler := NewAuthHandler(db, tracer, jwtSecret)
	
	// Test data
	username := "nonexistentuser"
	password := "password123"
	
	// Mock expectations
	mocks.ExpectGetUserNotFound(mock, username)
	
	// Call the method
	req := &pb.LoginRequest{
		Username: username,
		Password: password,
	}
	
	resp, err := handler.Login(context.Background(), req)
	
	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	
	// Check error code and message
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "user not found")
	
	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_InvalidPassword(t *testing.T) {
	// Setup
	db, mock, err := mocks.NewDBMock()
	assert.NoError(t, err)
	defer db.Close()

	tracer := mocks.NewMockTracer()
	jwtSecret := []byte("test-secret")
	
	handler := NewAuthHandler(db, tracer, jwtSecret)
	
	// Test data
	username := "testuser"
	password := "password123"
	wrongPassword := "wrongpassword"
	userID := int32(1)
	
	// Hash the password for storage
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	
	// Mock expectations
	mocks.ExpectGetUser(mock, username, userID, string(hashedPassword))
	
	// Call the method with wrong password
	req := &pb.LoginRequest{
		Username: username,
		Password: wrongPassword,
	}
	
	resp, err := handler.Login(context.Background(), req)
	
	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	
	// Check error code and message
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "invalid password")
	
	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHealth(t *testing.T) {
	// Setup
	db, _, err := mocks.NewDBMock()
	assert.NoError(t, err)
	defer db.Close()

	tracer := mocks.NewMockTracer()
	jwtSecret := []byte("test-secret")
	
	handler := NewAuthHandler(db, tracer, jwtSecret)
	
	// Call the method
	req := &pb.HealthRequest{}
	resp, err := handler.Health(context.Background(), req)
	
	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "ok", resp.Status)
} 