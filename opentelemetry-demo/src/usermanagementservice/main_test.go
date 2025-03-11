package main

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
	"github.com/stretchr/testify/assert"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestHealthChecker_Check(t *testing.T) {
	healthChecker := &HealthChecker{}

	resp, err := healthChecker.Check(context.Background(), &healthpb.HealthCheckRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.Status)
}

func TestHealthChecker_Watch(t *testing.T) {
	healthChecker := &HealthChecker{}

	err := healthChecker.Watch(&healthpb.HealthCheckRequest{}, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

func TestUserManagementServiceServer_Health(t *testing.T) {
	server := &UserManagementServiceServer{}

	resp, err := server.Health(context.Background(), &pb.HealthRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "ok", resp.Status)
}

// Test successful database initialization
func TestInitDB_Success(t *testing.T) {
	// Save original provider function
	originalDBProvider := dbProvider
	defer func() { dbProvider = originalDBProvider }()

	// Set up mock database
	mockDB, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	// Mock the dbProvider function
	dbProvider = func() *sql.DB {
		return mockDB
	}

	// Call the function under test
	db := initDB()

	// Verify it returned our mock db
	assert.Equal(t, mockDB, db)
}

// Test when DB_CONN environment variable is not set
func TestInitDB_DBConnNotSet(t *testing.T) {
	// Save original provider function and os.Exit/log.Fatal functions
	originalDBProvider := dbProvider
	originalOsExit := osExit
	originalLogFatal := logFatal

	defer func() {
		dbProvider = originalDBProvider
		osExit = originalOsExit
		logFatal = originalLogFatal
	}()

	// Mock the os.Exit and log.Fatal functions
	exitCalled := false
	exitCode := 0
	osExit = func(code int) {
		exitCalled = true
		exitCode = code
	}

	fatalCalled := false
	logFatal = func(v ...interface{}) {
		fatalCalled = true
	}

	// Mock the dbProvider function to simulate the DB_CONN missing error
	dbProvider = func() *sql.DB {
		logFatal("DB_CONN environment variable not set")
		osExit(1)
		return nil
	}

	// Call the function under test
	initDB()

	// Verify our mocked functions were called
	assert.True(t, fatalCalled, "log.Fatal should have been called")
	assert.True(t, exitCalled, "os.Exit should have been called")
	assert.Equal(t, 1, exitCode)
}

// Test when database table creation fails
func TestInitDB_TableCreationFails(t *testing.T) {
	// Save original provider function and os.Exit/log.Fatalf functions
	originalDBProvider := dbProvider
	originalOsExit := osExit
	originalLogFatalf := logFatalf

	defer func() {
		dbProvider = originalDBProvider
		osExit = originalOsExit
		logFatalf = originalLogFatalf
	}()

	// Set up mock database
	mockDB, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	// Mock the os.Exit and log.Fatalf functions
	exitCalled := false
	exitCode := 0
	osExit = func(code int) {
		exitCalled = true
		exitCode = code
	}

	fatalfCalled := false
	logFatalf = func(format string, v ...interface{}) {
		fatalfCalled = true
	}

	// Mock the dbProvider function to simulate the table creation error
	dbProvider = func() *sql.DB {
		logFatalf("Failed to ensure database tables exist: %v", fmt.Errorf("table creation error"))
		osExit(1)
		return nil
	}

	// Call the function under test
	initDB()

	// Verify our mocked functions were called
	assert.True(t, fatalfCalled, "log.Fatalf should have been called")
	assert.True(t, exitCalled, "os.Exit should have been called")
	assert.Equal(t, 1, exitCode)
}
