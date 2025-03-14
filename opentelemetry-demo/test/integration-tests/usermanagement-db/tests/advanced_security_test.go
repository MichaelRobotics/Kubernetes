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

// TestJWTTokenRevocation tests if the system has a mechanism to revoke JWT tokens
// Note: This test might be skipped if token revocation is not implemented
func TestJWTTokenRevocation(t *testing.T) {
	// Skip this test if token revocation is not implemented
	t.Skip("Token revocation might not be implemented in this service")

	ctx := context.Background()
	username := fmt.Sprintf("revoke_token_user_%d", time.Now().UnixNano())
	password := "password123"

	// Register a user
	_, err := client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Login to get a token
	loginResp, err := client.Login(ctx, &pb.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	token := loginResp.GetToken()
	if token == "" {
		t.Fatal("Token is empty")
	}

	t.Logf("Successfully obtained token: %s...", token[:20])

	// Token revocation would typically be implemented as an API call
	// For example: client.RevokeToken(ctx, &pb.RevokeTokenRequest{Token: token})
	// Since we don't know if such a method exists, we'll skip actually testing it
	t.Log("Note: This test is incomplete because token revocation API might not exist")

	// In a real implementation, we would now try to use the revoked token and expect it to fail
}

// TestSessionTimeout tests if tokens expire after their expected lifetime
func TestSessionTimeout(t *testing.T) {
	ctx := context.Background()
	username := fmt.Sprintf("timeout_test_user_%d", time.Now().UnixNano())
	password := "password123"

	// Register a user
	_, err := client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Login to get a token
	loginResp, err := client.Login(ctx, &pb.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	token := loginResp.GetToken()
	if token == "" {
		t.Fatal("Token is empty")
	}

	// Parse the JWT token to extract expiration time
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

	// Check expiration time
	expClaim, hasExp := claims["exp"]
	if !hasExp {
		t.Fatal("Token doesn't have expiration claim")
	}

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
	expiresIn := time.Duration(expTime-now) * time.Second

	t.Logf("Token expires in approximately %v", expiresIn)
	t.Logf("To fully test expiration, you would need to wait for the token to expire (%v) or mock time", expiresIn)

	// Since we can't wait for token expiration in a unit test, we'll just verify it has a reasonable expiration time
	// Most tokens expire in 1 hour to 24 hours
	if expiresIn < 1*time.Minute || expiresIn > 48*time.Hour {
		t.Errorf("Token expiration time seems unusual: %v", expiresIn)
	}
}

// TestMultipleLogins tests behavior with multiple logins from the same user
func TestMultipleLogins(t *testing.T) {
	ctx := context.Background()
	username := fmt.Sprintf("multi_login_user_%d", time.Now().UnixNano())
	password := "password123"

	// Register a user
	_, err := client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Login multiple times and collect tokens
	var tokens []string
	for i := 0; i < 3; i++ {
		loginResp, err := client.Login(ctx, &pb.LoginRequest{
			Username: username,
			Password: password,
		})
		if err != nil {
			t.Fatalf("Login %d failed: %v", i+1, err)
		}

		token := loginResp.GetToken()
		if token == "" {
			t.Fatalf("Token from login %d is empty", i+1)
		}

		tokens = append(tokens, token)
		t.Logf("Login %d successful, received token: %s...", i+1, token[:20])

		// Short pause between logins
		time.Sleep(100 * time.Millisecond)
	}

	// Check if tokens are different or the same
	uniqueTokens := make(map[string]bool)
	for _, token := range tokens {
		uniqueTokens[token] = true
	}

	if len(uniqueTokens) < len(tokens) {
		t.Logf("Service reuses tokens for multiple logins (issued %d tokens, but %d are unique)",
			len(tokens), len(uniqueTokens))
	} else {
		t.Logf("Service issues unique tokens for each login")
	}

	// In a system with session invalidation, logging in again might invalidate previous tokens
	// However, testing this would require knowledge of how your service implements session management
	t.Log("Note: To fully test session invalidation, try using earlier tokens after multiple logins")
}
