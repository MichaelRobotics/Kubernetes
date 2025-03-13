package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestMalformedRequests verifies that the service properly handles malformed requests
func TestMalformedRequests(t *testing.T) {
	ctx := context.Background()

	// Test with an extremely long username (beyond any reasonable limit)
	veryLongUsername := ""
	for i := 0; i < 1000; i++ {
		veryLongUsername += "a"
	}

	_, err := client.Register(ctx, &pb.RegisterRequest{
		Username: veryLongUsername,
		Password: "validpassword123",
	})

	if err == nil {
		t.Error("Expected error when registering with extremely long username, but got nil")
	} else {
		t.Logf("Service correctly rejected extremely long username: %v", err)
	}

	// Test with potentially problematic characters in username
	problematicCharacters := []string{
		"<script>alert('xss')</script>",
		"';DROP TABLE users;--",
		"null",
		"undefined",
		"true",
		"false",
		"{\"key\":\"value\"}",
	}

	for _, problematic := range problematicCharacters {
		_, err := client.Register(ctx, &pb.RegisterRequest{
			Username: problematic,
			Password: "validpassword123",
		})

		if err == nil {
			// If registration succeeded, verify that the problematic username can be used to login
			loginResp, loginErr := client.Login(ctx, &pb.LoginRequest{
				Username: problematic,
				Password: "validpassword123",
			})

			if loginErr != nil {
				t.Errorf("Service accepted problematic username '%s' for registration but login failed: %v", problematic, loginErr)
			} else {
				t.Logf("Service accepted problematic username '%s' and login succeeded with user ID: %d", problematic, loginResp.GetUserId())
			}
		} else {
			t.Logf("Service rejected problematic username '%s': %v", problematic, err)
		}
	}
}

// TestMissingRequiredFields verifies that the service properly handles requests with missing required fields
func TestMissingRequiredFields(t *testing.T) {
	ctx := context.Background()

	// Test register with missing username
	_, err := client.Register(ctx, &pb.RegisterRequest{
		Username: "",
		Password: "validpassword123",
	})

	if err == nil {
		t.Error("Expected error when registering with empty username, but got nil")
	} else {
		errStatus, ok := status.FromError(err)
		if !ok {
			t.Errorf("Expected gRPC status error, got: %v", err)
		} else if errStatus.Code() != codes.InvalidArgument {
			t.Errorf("Expected InvalidArgument error for empty username, got %s: %s", errStatus.Code(), errStatus.Message())
		} else {
			t.Logf("Service correctly rejected empty username with InvalidArgument: %v", err)
		}
	}

	// Test register with missing password
	username := fmt.Sprintf("test_user_%d", time.Now().UnixNano())
	_, err = client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: "",
	})

	if err == nil {
		t.Error("Expected error when registering with empty password, but got nil")
	} else {
		errStatus, ok := status.FromError(err)
		if !ok {
			t.Errorf("Expected gRPC status error, got: %v", err)
		} else if errStatus.Code() != codes.InvalidArgument {
			t.Errorf("Expected InvalidArgument error for empty password, got %s: %s", errStatus.Code(), errStatus.Message())
		} else {
			t.Logf("Service correctly rejected empty password with InvalidArgument: %v", err)
		}
	}

	// Test login with missing username
	_, err = client.Login(ctx, &pb.LoginRequest{
		Username: "",
		Password: "validpassword123",
	})

	if err == nil {
		t.Error("Expected error when logging in with empty username, but got nil")
	} else {
		t.Logf("Service correctly rejected login with empty username: %v", err)
	}

	// Test login with missing password
	_, err = client.Login(ctx, &pb.LoginRequest{
		Username: "some_user",
		Password: "",
	})

	if err == nil {
		t.Error("Expected error when logging in with empty password, but got nil")
	} else {
		t.Logf("Service correctly rejected login with empty password: %v", err)
	}
}
