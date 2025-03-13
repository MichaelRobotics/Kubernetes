package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
)

// TestUnicodeCharacters verifies that the system properly handles Unicode characters in usernames and passwords
func TestUnicodeCharacters(t *testing.T) {
	ctx := context.Background()

	// Test various Unicode characters in usernames
	unicodeUsernames := []string{
		"résumé",     // Accented Latin characters
		"用户名",        // Chinese characters
		"Имя",        // Cyrillic
		"사용자",        // Korean
		"उपयोगकर्ता", // Devanagari
		"משתמש",      // Hebrew
		"مستخدم",     // Arabic
		"😀👍🔥",        // Emoji
	}

	unicodePasswords := []string{
		"pass_résumé_word", // Accented Latin characters
		"密码用户名123",         // Chinese characters
		"пароль123",        // Cyrillic
		"비밀번호123",          // Korean
		"पासवर्ड123",       // Devanagari
		"סיסמה123",         // Hebrew
		"كلمة المرور123",   // Arabic
		"password_😀👍🔥",     // Emoji
	}

	for i, username := range unicodeUsernames {
		// Use a corresponding password from the list or fall back to a default
		password := "unicodepass123"
		if i < len(unicodePasswords) {
			password = unicodePasswords[i]
		}

		// Attempt to register
		registerResp, err := client.Register(ctx, &pb.RegisterRequest{
			Username: username,
			Password: password,
		})

		if err != nil {
			t.Logf("Registration failed for Unicode username '%s': %v", username, err)
		} else {
			t.Logf("Successfully registered user with Unicode username '%s', user ID: %d",
				username, registerResp.GetUserId())

			// Try to login with the same credentials
			loginResp, err := client.Login(ctx, &pb.LoginRequest{
				Username: username,
				Password: password,
			})

			if err != nil {
				t.Errorf("Login failed for Unicode username '%s': %v", username, err)
			} else {
				t.Logf("Successfully logged in with Unicode username '%s', token: %s...",
					username, loginResp.GetToken()[:20])
			}
		}
	}
}

// TestMaxConnections tests the system's behavior under maximum concurrent connections
func TestMaxConnections(t *testing.T) {
	ctx := context.Background()

	// Large number of concurrent connections
	// Note: Adjust this number based on your system's expected capacity
	numConnections := 100

	var wg sync.WaitGroup
	results := make(chan error, numConnections)

	t.Logf("Testing with %d concurrent connections", numConnections)
	startTime := time.Now()

	// Launch many concurrent operations
	for i := 0; i < numConnections; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			// Generate a unique username
			username := fmt.Sprintf("max_conn_user_%d_%d", time.Now().UnixNano(), idx)

			// Try to register a user
			_, err := client.Register(ctx, &pb.RegisterRequest{
				Username: username,
				Password: "maxconnpass123",
			})

			if err != nil {
				results <- fmt.Errorf("connection %d failed: %v", idx, err)
			} else {
				results <- nil
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)

	// Analyze results
	totalTime := time.Since(startTime)
	successCount := 0
	failureCount := 0

	// Categorize errors
	errorTypes := make(map[string]int)

	for err := range results {
		if err == nil {
			successCount++
		} else {
			failureCount++
			// Count error types to detect connection limits
			errorStr := err.Error()
			errorTypes[errorStr]++

			// Only log a few examples to avoid flooding the output
			if failureCount <= 5 {
				t.Logf("Connection error: %v", err)
			}
		}
	}

	t.Logf("Max connections test completed in %v", totalTime)
	t.Logf("Successful connections: %d/%d", successCount, numConnections)
	t.Logf("Failed connections: %d/%d", failureCount, numConnections)

	if failureCount > 0 {
		t.Logf("Error types encountered:")
		for errMsg, count := range errorTypes {
			t.Logf("  - %s: %d occurrences", errMsg, count)
		}
	}

	// Check if we hit connection limits
	if failureCount > numConnections/2 {
		t.Logf("WARNING: More than half of connections failed, possibly hit connection limits")
	}
}

// TestDatabaseRecovery tests if the system can recover from database connection issues
// Note: This is a more advanced test that requires manipulating the database connection
func TestDatabaseRecovery(t *testing.T) {
	// Skip this test if we don't have access to manipulate the database container
	t.Skip("Database recovery test requires direct access to database container controls")

	/*
		// This test would typically:
		// 1. Register a user successfully
		// 2. Force the database to become unavailable (e.g., stop postgres container)
		// 3. Verify that operations fail gracefully
		// 4. Restore the database (e.g., start postgres container)
		// 5. Verify that operations succeed again after recovery

		ctx := context.Background()
		username := fmt.Sprintf("recovery_test_user_%d", time.Now().UnixNano())

		// 1. Register a user
		_, err := client.Register(ctx, &pb.RegisterRequest{
			Username: username,
			Password: "recoverypass123",
		})
		if err != nil {
			t.Fatalf("Failed to register user before database disruption: %v", err)
		}

		// 2. Stop the database
		t.Log("Stopping database container...")
		cmd := exec.Command("docker", "stop", "ums-db")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to stop database container: %v", err)
		}

		// Wait for database to stop
		time.Sleep(2 * time.Second)

		// 3. Verify operations fail gracefully
		_, err = client.Login(ctx, &pb.LoginRequest{
			Username: username,
			Password: "recoverypass123",
		})
		if err == nil {
			t.Error("Expected login to fail when database is down, but it succeeded")
		} else {
			t.Logf("Login correctly failed when database is down: %v", err)
		}

		// 4. Restart the database
		t.Log("Restarting database container...")
		cmd = exec.Command("docker", "start", "ums-db")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to restart database container: %v", err)
		}

		// Wait for database to recover
		time.Sleep(10 * time.Second)

		// 5. Verify operations succeed again
		_, err = client.Login(ctx, &pb.LoginRequest{
			Username: username,
			Password: "recoverypass123",
		})
		if err != nil {
			t.Errorf("Login failed after database recovery: %v", err)
		} else {
			t.Log("Successfully logged in after database recovery")
		}
	*/
}
