package main

import (
	"context"
	"strings"
	"testing"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
)

// TestUsernameValidation verifies that the system enforces username requirements
func TestUsernameValidation(t *testing.T) {
	ctx := context.Background()

	// Test empty username
	_, err := client.Register(ctx, &pb.RegisterRequest{
		Username: "",
		Password: "validpassword12345",
	})
	if err == nil {
		t.Error("Expected error when registering with empty username")
	} else {
		t.Logf("Correctly rejected empty username: %v", err)
	}

	// Test very long username
	veryLongUsername := strings.Repeat("a", 100)
	_, err = client.Register(ctx, &pb.RegisterRequest{
		Username: veryLongUsername,
		Password: "validpassword12345",
	})
	if err == nil {
		t.Log("Service accepts 100-character usernames")
	} else {
		t.Logf("Service rejected long username: %v", err)
	}

	// Test special characters
	specialChars := []string{"user@name", "user name", "user#name", "user;name"}
	for _, username := range specialChars {
		_, err = client.Register(ctx, &pb.RegisterRequest{
			Username: username,
			Password: "validpassword12345",
		})
		if err == nil {
			t.Logf("Service accepts special characters in username: %s", username)
		} else {
			t.Logf("Service rejected special characters in username: %s - %v", username, err)
		}
	}
}
