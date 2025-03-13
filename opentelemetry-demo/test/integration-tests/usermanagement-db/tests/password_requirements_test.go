package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
)

// TestPasswordRequirements verifies that the system enforces password requirements if they exist
func TestPasswordRequirements(t *testing.T) {
	ctx := context.Background()

	// Test empty password
	username := fmt.Sprintf("emptypass_user_%d", time.Now().UnixNano())
	_, err := client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: "",
	})

	// Check if the service validates empty passwords
	// If it does, we should get an error
	if err == nil {
		t.Log("Warning: Service allows empty passwords")
	} else {
		t.Logf("Service correctly rejected empty password: %v", err)
	}

	// Test very short password
	username = fmt.Sprintf("shortpass_user_%d", time.Now().UnixNano())
	_, err = client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: "a",
	})

	// If the service has min length requirements, this should fail
	if err == nil {
		t.Log("Warning: Service allows single-character passwords")
	} else {
		t.Logf("Service correctly rejected very short password: %v", err)
	}

	// Test very long password
	username = fmt.Sprintf("longpass_user_%d", time.Now().UnixNano())
	veryLongPassword := ""
	for i := 0; i < 1000; i++ {
		veryLongPassword += "a"
	}

	_, err = client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: veryLongPassword,
	})

	// If the service has max length requirements, this might fail
	if err == nil {
		t.Logf("Service accepted a 1000-character password")
	} else {
		t.Logf("Service rejected very long password: %v", err)
	}
}
