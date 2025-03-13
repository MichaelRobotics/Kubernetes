package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
)

// TestComprehensiveUserManagement runs a comprehensive series of operations to simulate
// real-world usage patterns and verify the system behaves correctly
func TestComprehensiveUserManagement(t *testing.T) {
	ctx := context.Background()

	// -- Part 1: Basic User Lifecycle --

	// Generate unique user
	username := fmt.Sprintf("comprehensive_user_%d", time.Now().UnixNano())
	password := "Compl3xP@ssw0rd!"

	t.Logf("Testing comprehensive user management flow for user: %s", username)

	// Registration
	registerResp, err := client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}
	userID := registerResp.GetUserId()
	t.Logf("User registered successfully with ID: %d", userID)

	// Login with correct credentials
	loginResp, err := client.Login(ctx, &pb.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Fatalf("Login failed with correct credentials: %v", err)
	}
	token := loginResp.GetToken()
	t.Logf("Login successful, received token: %s...", token[:20])

	// -- Part 2: Error Handling --

	// Try registering the same user again (should fail)
	_, err = client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: "AnotherPassword123",
	})
	if err == nil {
		t.Error("Expected error when registering duplicate username, but got nil")
	} else {
		t.Logf("Correctly rejected duplicate registration: %v", err)
	}

	// Try logging in with wrong password
	_, err = client.Login(ctx, &pb.LoginRequest{
		Username: username,
		Password: "wrongpassword",
	})
	if err == nil {
		t.Error("Expected error when logging in with wrong password, but got nil")
	} else {
		t.Logf("Correctly rejected login with wrong password: %v", err)
	}

	// -- Part 3: Persistence Testing --

	// Verify user can log in multiple times
	for i := 0; i < 3; i++ {
		_, err = client.Login(ctx, &pb.LoginRequest{
			Username: username,
			Password: password,
		})
		if err != nil {
			t.Errorf("Login attempt %d failed unexpectedly: %v", i+1, err)
		}
		// Brief pause between logins
		time.Sleep(100 * time.Millisecond)
	}
	t.Log("Multiple logins successful")

	// -- Part 4: Boundary Testing --

	// Test with almost-empty username (minimal length)
	// Most systems require at least 3 characters
	minUsername := fmt.Sprintf("min%d", time.Now().UnixNano())
	minRegisterResp, err := client.Register(ctx, &pb.RegisterRequest{
		Username: minUsername,
		Password: password,
	})
	if err != nil {
		t.Logf("Registration with minimal username failed: %v", err)
	} else {
		t.Logf("Registration with minimal username succeeded, user ID: %d", minRegisterResp.GetUserId())

		// Try logging in with minimal username
		_, err = client.Login(ctx, &pb.LoginRequest{
			Username: minUsername,
			Password: password,
		})
		if err != nil {
			t.Errorf("Login with minimal username failed: %v", err)
		} else {
			t.Log("Login with minimal username succeeded")
		}
	}

	// -- Part 5: Special Characters Testing --

	// Test with username containing special characters
	specialUsername := fmt.Sprintf("special_user-%d@test", time.Now().UnixNano())
	specialRegisterResp, err := client.Register(ctx, &pb.RegisterRequest{
		Username: specialUsername,
		Password: password,
	})
	if err != nil {
		t.Logf("Registration with special characters in username failed: %v", err)
	} else {
		t.Logf("Registration with special characters in username succeeded, user ID: %d",
			specialRegisterResp.GetUserId())

		// Try logging in with special username
		_, err = client.Login(ctx, &pb.LoginRequest{
			Username: specialUsername,
			Password: password,
		})
		if err != nil {
			t.Errorf("Login with special characters in username failed: %v", err)
		} else {
			t.Log("Login with special characters in username succeeded")
		}
	}

	// -- Final Summary --
	t.Log("Comprehensive test completed successfully")
}
