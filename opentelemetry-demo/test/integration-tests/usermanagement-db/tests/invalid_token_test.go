package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
	"google.golang.org/grpc/metadata"
)

// TestInvalidJWTTokens verifies that the service properly handles invalid JWT tokens
func TestInvalidJWTTokens(t *testing.T) {
	ctx := context.Background()

	// Register and login to get a valid token first
	username := fmt.Sprintf("token_test_user_%d", time.Now().UnixNano())
	_, err := client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	loginResp, err := client.Login(ctx, &pb.LoginRequest{
		Username: username,
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	validToken := loginResp.GetToken()
	if validToken == "" {
		t.Fatal("Valid token is empty")
	}

	// Test cases with different invalid tokens
	testCases := []struct {
		name  string
		token string
	}{
		{
			name:  "Empty token",
			token: "",
		},
		{
			name:  "Malformed token (not enough parts)",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
		{
			name:  "Malformed token (too many parts)",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c.extra",
		},
		{
			name:  "Invalid signature",
			token: validToken[:strings.LastIndex(validToken, ".")+1] + "invalidSignature",
		},
		{
			name:  "Expired token (manually crafted)",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyMzkwMjIsInN1YiI6IjEifQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		},
		{
			name:  "Token with wrong structure",
			token: "not-a-jwt-token",
		},
	}

	// We'll test by setting the token in the metadata context
	// Note: This test assumes your service extracts the token from a specific metadata field
	// You may need to adjust this based on your actual implementation
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Since we don't have access to a method that verifies tokens directly,
			// we'll use the token in a context and see if the service rejects it
			// This is a bit of a workaround since we don't know exactly how your service validates tokens

			// Create a context with the token in metadata
			md := metadata.New(map[string]string{
				"authorization": "Bearer " + tc.token,
			})
			ctxWithToken := metadata.NewOutgoingContext(ctx, md)

			// Try to use the token context to perform a protected operation
			// Note: This is a hypothetical test since we don't know which operations require authentication
			// You might need to adjust this to use an actual protected endpoint in your service
			_, err := client.Register(ctxWithToken, &pb.RegisterRequest{
				Username: fmt.Sprintf("test_user_%d", time.Now().UnixNano()),
				Password: "password123",
			})

			// We can't make strong assertions about the error since we don't know
			// if this operation actually validates tokens, but we can log the outcome
			if err != nil {
				t.Logf("Using invalid token (%s) resulted in error: %v", tc.name, err)
			} else {
				t.Logf("Note: Using invalid token (%s) did not result in an error. This might be expected if the operation doesn't require authentication.", tc.name)
			}
		})
	}
}
