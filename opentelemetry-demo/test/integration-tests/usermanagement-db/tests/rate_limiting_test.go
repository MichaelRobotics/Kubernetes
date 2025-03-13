package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
)

// TestRateLimiting verifies if the system has rate limiting protection
func TestRateLimiting(t *testing.T) {
	ctx := context.Background()

	// Make a large number of requests in a short time
	requestCount := 50
	results := make(chan string, requestCount)
	var wg sync.WaitGroup

	// Create a username that we'll repeatedly try to authenticate
	username := fmt.Sprintf("ratelimit_user_%d", time.Now().UnixNano())

	// Register the user first
	_, err := client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Failed to register user for rate limit test: %v", err)
	}

	startTime := time.Now()

	// Make many login requests rapidly
	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			_, err := client.Login(ctx, &pb.LoginRequest{
				Username: username,
				Password: "password123",
			})

			if err != nil {
				// Check if it's a rate limiting error (typically contains words like "rate", "limit", "too many")
				errMsg := err.Error()
				if strings.Contains(strings.ToLower(errMsg), "rate") ||
					strings.Contains(strings.ToLower(errMsg), "limit") ||
					strings.Contains(strings.ToLower(errMsg), "too many") {
					results <- fmt.Sprintf("rate limit detected: %v", err)
				} else {
					results <- fmt.Sprintf("error: %v", err)
				}
			} else {
				results <- "success"
			}
		}(i)
	}

	wg.Wait()
	close(results)

	// Analyze results
	totalTime := time.Since(startTime)
	successCount := 0
	errorCount := 0
	rateLimitCount := 0

	for result := range results {
		if result == "success" {
			successCount++
		} else if strings.HasPrefix(result, "rate limit detected:") {
			rateLimitCount++
		} else {
			errorCount++
		}
	}

	t.Logf("Rate limiting test results (total time: %v):", totalTime)
	t.Logf("  - Successful requests: %d", successCount)
	t.Logf("  - Rate-limited requests: %d", rateLimitCount)
	t.Logf("  - Other errors: %d", errorCount)

	if rateLimitCount > 0 {
		t.Logf("Rate limiting detected: %d requests were rate-limited", rateLimitCount)
	} else {
		t.Logf("No rate limiting detected for %d concurrent requests in %v", requestCount, totalTime)
		if totalTime < 3*time.Second {
			t.Logf("Warning: All %d requests completed in %v without rate limiting", requestCount, totalTime)
		}
	}
}
