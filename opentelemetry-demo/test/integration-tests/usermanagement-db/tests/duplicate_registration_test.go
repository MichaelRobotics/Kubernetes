package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
)

// TestDuplicateRegistration verifies that the system prevents users from registering with an existing username
func TestDuplicateRegistration(t *testing.T) {
	ctx := context.Background()
	username := fmt.Sprintf("duplicate_user_%d", time.Now().UnixNano())

	// Register the first user
	registerResp1, err := client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	t.Logf("Successfully registered first user with ID: %d", registerResp1.GetUserId())

	// Try to register the same username again
	_, err = client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: "differentpassword",
	})

	// Should fail with an appropriate error
	if err == nil {
		t.Fatal("Expected error for duplicate username, but got nil")
	}

	t.Logf("Correctly received error when trying to register duplicate username: %v", err)
}
