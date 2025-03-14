package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
)

// TestSQLInjectionProtection tests if the service properly sanitizes inputs to prevent SQL injection attacks
func TestSQLInjectionProtection(t *testing.T) {
	ctx := context.Background()

	// Common SQL injection patterns for username field
	sqlInjectionPatterns := []string{
		"' OR '1'='1",
		"admin'; --",
		"username' OR 1=1; --",
		"'; DROP TABLE users; --",
		"' UNION SELECT username, password FROM users; --",
		"1' or '1' = '1",
		"1' or 1=1--",
		"' or 1=1--",
		"' or ''='",
		"' or 1 --'",
		"' or 1/*",
		"' or '1'='1'--",
	}

	// Test SQL injection through username field during login
	for _, pattern := range sqlInjectionPatterns {
		// Try login with SQL injection pattern
		_, err := client.Login(ctx, &pb.LoginRequest{
			Username: pattern,
			Password: "anypassword",
		})

		if err == nil {
			t.Errorf("Potential SQL injection vulnerability: login successful with pattern '%s'", pattern)
		} else {
			t.Logf("SQL injection pattern rejected in login: %s - %v", pattern, err)
		}
	}

	// Test SQL injection through username field during registration
	for idx, pattern := range sqlInjectionPatterns {
		uniqueSuffix := fmt.Sprintf("_%d_%d", time.Now().UnixNano(), idx)
		injectionUsername := pattern + uniqueSuffix

		_, err := client.Register(ctx, &pb.RegisterRequest{
			Username: injectionUsername,
			Password: "validpassword12345",
		})

		if err == nil {
			// Try to login with the same username to see if injection worked
			_, loginErr := client.Login(ctx, &pb.LoginRequest{
				Username: injectionUsername,
				Password: "validpassword12345",
			})

			if loginErr == nil {
				t.Logf("Registration allowed potentially dangerous username: %s (and login succeeded)", injectionUsername)
			} else {
				t.Logf("Registration allowed potentially dangerous username: %s (but login failed: %v)", injectionUsername, loginErr)
			}
		} else {
			t.Logf("Registration rejected potentially dangerous username: %s - %v", injectionUsername, err)
		}
	}

	// Test SQL injection in password field
	// This is less likely to be vulnerable but worth testing
	cleanUsername := fmt.Sprintf("clean_user_%d", time.Now().UnixNano())

	_, err := client.Register(ctx, &pb.RegisterRequest{
		Username: cleanUsername,
		Password: "validpassword", // Register with valid password first
	})
	if err != nil {
		t.Fatalf("Failed to register clean user: %v", err)
	}

	// Try login with SQL injection in password
	for _, pattern := range sqlInjectionPatterns {
		_, err := client.Login(ctx, &pb.LoginRequest{
			Username: cleanUsername,
			Password: pattern,
		})

		if err == nil {
			t.Errorf("Potential SQL injection vulnerability: login successful with password pattern '%s'", pattern)
		} else {
			t.Logf("SQL injection pattern in password correctly rejected: %s - %v", pattern, err)
		}
	}
}
