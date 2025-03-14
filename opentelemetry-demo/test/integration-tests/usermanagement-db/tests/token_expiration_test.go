package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
)

// TestTokenExpiration verifies that tokens have an expiration timestamp
func TestTokenExpiration(t *testing.T) {
	ctx := context.Background()
	username := fmt.Sprintf("expiration_user_%d", time.Now().UnixNano())

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

	// Parse the token to examine expiration
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("Token doesn't have expected format: %s", token)
	}

	// Decode the payload (second part)
	// Note: JWT padding is removed, so we need to add it back for base64 decoding
	payload := parts[1]
	if len(payload)%4 != 0 {
		payload += strings.Repeat("=", 4-len(payload)%4)
	}

	decodedBytes, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		t.Fatalf("Failed to decode payload: %v", err)
	}

	// Extract claims
	var claims map[string]interface{}
	if err := json.Unmarshal(decodedBytes, &claims); err != nil {
		t.Fatalf("Failed to parse claims: %v", err)
	}

	// Check for expiration claim
	expClaim, hasExp := claims["exp"]
	if !hasExp {
		t.Fatal("Token doesn't have expiration claim")
	}

	// Verify expiration is in the future
	var expTime float64
	switch exp := expClaim.(type) {
	case float64:
		expTime = exp
	case json.Number:
		expTime, _ = exp.Float64()
	default:
		t.Fatalf("Unexpected type for exp claim: %T", expClaim)
	}

	now := float64(time.Now().Unix())
	if expTime <= now {
		t.Errorf("Token expiration time %f is not in the future (current time: %f)", expTime, now)
	} else {
		// Calculate how long until expiration
		validFor := time.Duration(expTime-now) * time.Second
		t.Logf("Token is valid for approximately %v", validFor)
	}
}
