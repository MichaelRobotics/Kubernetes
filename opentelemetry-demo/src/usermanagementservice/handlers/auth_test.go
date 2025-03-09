package handlers_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
	"github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/handlers"
	"github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/tests/mocks"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace/noop"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func setupAuthHandler() (*handlers.AuthHandler, sqlmock.Sqlmock, *sql.DB) {
	db, mock, _ := mocks.MockDB()
	tracer := noop.NewTracerProvider().Tracer("test-tracer")
	jwtSecret := []byte("test-secret")
	handler := handlers.NewAuthHandler(db, tracer, jwtSecret)
	return handler, mock, db
}

// Register endpoint tests

func TestRegister_Success(t *testing.T) {
	handler, mock, db := setupAuthHandler()
	defer db.Close()

	req := &pb.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	}

	// First query - check if username exists
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(req.Username).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	// Second query - insert user and return ID
	mock.ExpectQuery("INSERT INTO users").
		WithArgs(req.Username, sqlmock.AnyArg()). // Can't predict exact hash
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	resp, err := handler.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(1), resp.UserId)
	assert.Equal(t, req.Username, resp.Username)
	assert.Equal(t, "User registered successfully", resp.Message)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegister_UsernameTooShort(t *testing.T) {
	handler, mock, db := setupAuthHandler()
	defer db.Close()

	req := &pb.RegisterRequest{
		Username: "ab", // Too short
		Password: "password123",
	}

	resp, err := handler.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "username must be at least 3 characters")
	assert.NoError(t, mock.ExpectationsWereMet()) // No DB calls
}

func TestRegister_PasswordTooShort(t *testing.T) {
	handler, mock, db := setupAuthHandler()
	defer db.Close()

	req := &pb.RegisterRequest{
		Username: "testuser",
		Password: "short", // Too short
	}

	resp, err := handler.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "password must be at least 8 characters")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegister_UsernameExists(t *testing.T) {
	handler, mock, db := setupAuthHandler()
	defer db.Close()

	req := &pb.RegisterRequest{
		Username: "existinguser",
		Password: "password123",
	}

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(req.Username).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	resp, err := handler.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, st.Code())
	assert.Contains(t, st.Message(), "username already exists")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegister_DatabaseError(t *testing.T) {
	handler, mock, db := setupAuthHandler()
	defer db.Close()

	req := &pb.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	}

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(req.Username).
		WillReturnError(sql.ErrConnDone)

	resp, err := handler.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Login endpoint tests

func TestLogin_Success(t *testing.T) {
	handler, mock, db := setupAuthHandler()
	defer db.Close()

	req := &pb.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	// Create valid hash for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	mock.ExpectQuery("SELECT id, username, password_hash FROM users").
		WithArgs(req.Username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash"}).
			AddRow(1, req.Username, string(hashedPassword)))

	resp, err := handler.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(1), resp.UserId)
	assert.NotEmpty(t, resp.Token)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_UserNotFound(t *testing.T) {
	handler, mock, db := setupAuthHandler()
	defer db.Close()

	req := &pb.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}

	mock.ExpectQuery("SELECT id, username, password_hash FROM users").
		WithArgs(req.Username).
		WillReturnError(sql.ErrNoRows)

	resp, err := handler.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "user not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_InvalidPassword(t *testing.T) {
	handler, mock, db := setupAuthHandler()
	defer db.Close()

	req := &pb.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	// Hash for a different password
	correctPassword := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

	mock.ExpectQuery("SELECT id, username, password_hash FROM users").
		WithArgs(req.Username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash"}).
			AddRow(1, req.Username, string(hashedPassword)))

	resp, err := handler.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "invalid password")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Health endpoint test

func TestHealth(t *testing.T) {
	handler, _, db := setupAuthHandler()
	defer db.Close()

	resp, err := handler.Health(context.Background(), &pb.HealthRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "ok", resp.Status)
}
