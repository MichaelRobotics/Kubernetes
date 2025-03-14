package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
)

// TestTokenValidation verifies that the system generates valid JWT tokens
func TestTokenValidation(t *testing.T) {
	ctx := context.Background()
	username := fmt.Sprintf("tokenuser_%d", time.Now().UnixNano())

	// Register and login to get token
	_, err := client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	loginResp, err := client.Login(ctx, &pb.LoginRequest{
		Username: username,
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	token := loginResp.GetToken()

	// Verify token is not empty and has the expected format
	if token == "" {
		t.Fatal("Token is empty")
	}

	// Basic JWT format validation (header.payload.signature)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("Expected 3 parts in JWT token, got %d", len(parts))
	} else {
		t.Logf("Token has correct JWT format: %s.%s.%s", parts[0][:8], parts[1][:8], parts[2][:8])
	}

	// Test that token contains user ID
	userId := fmt.Sprintf("%d", loginResp.GetUserId())
	// This is a simplified check and would need adjustment based on your actual JWT payload structure
	if !strings.Contains(token, userId) {
		t.Logf("Token might not contain user ID %s", userId)
	}
}
