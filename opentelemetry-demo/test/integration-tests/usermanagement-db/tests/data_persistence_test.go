package main

import (
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
)

// TestDataPersistence verifies that user data persists after service restart
func TestDataPersistence(t *testing.T) {
	// The following code is adapted to the Docker environment
	ctx := context.Background()
	username := fmt.Sprintf("persistent_user_%d", time.Now().UnixNano())
	password := "persistentpass"

	// Register a user
	registerResp, err := client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}
	userId := registerResp.GetUserId()
	t.Logf("Registered user with ID: %d", userId)

	// Restart the user management service
	t.Log("Restarting user management service...")
	cmd := exec.Command("docker", "restart", "ums-service")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to restart service container: %v", err)
	}

	// Wait for service to restart
	time.Sleep(5 * time.Second)

	// Try to login with the same credentials after restart
	loginResp, err := client.Login(ctx, &pb.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Fatalf("Login after service restart failed: %v", err)
	}

	// Verify the user ID is preserved
	if loginResp.GetUserId() != userId {
		t.Errorf("User ID changed after restart: expected %d, got %d",
			userId, loginResp.GetUserId())
	} else {
		t.Logf("User ID successfully persisted after service restart: %d", userId)
	}
}
