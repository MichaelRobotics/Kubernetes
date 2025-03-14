package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
)

// TestInvalidLogin verifies that the system correctly handles invalid login attempts
func TestInvalidLogin(t *testing.T) {
	ctx := context.Background()

	// Test with non-existent user
	nonExistentUser := fmt.Sprintf("nonexistent_%d", time.Now().UnixNano())
	_, err := client.Login(ctx, &pb.LoginRequest{
		Username: nonExistentUser,
		Password: "somepassword",
	})
	if err == nil {
		t.Error("Expected error when logging in with non-existent user")
	} else {
		t.Logf("Correctly received error for non-existent user: %v", err)
	}

	// Register a user then try wrong password
	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	_, err = client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: "correctpassword",
	})
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	// Try login with wrong password
	_, err = client.Login(ctx, &pb.LoginRequest{
		Username: username,
		Password: "wrongpassword",
	})
	if err == nil {
		t.Error("Expected error when logging in with incorrect password")
	} else {
		t.Logf("Correctly received error for wrong password: %v", err)
	}
}
