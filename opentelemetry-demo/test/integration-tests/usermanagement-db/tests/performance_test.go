package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
)

// TestPerformance tests the response time of the service under load
func TestPerformance(t *testing.T) {
	ctx := context.Background()

	// Adjust these values based on your performance expectations
	userCount := 20
	maxAvgResponseTime := 500 * time.Millisecond // Generous threshold for testing

	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make(map[int]struct {
		registerTime time.Duration
		loginTime    time.Duration
		error        error
	})

	t.Logf("Starting performance test with %d concurrent users", userCount)
	startTime := time.Now()

	for i := 0; i < userCount; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			// Generate unique username
			username := fmt.Sprintf("perf_user_%d_%d", time.Now().UnixNano(), idx)

			// Structure to store results
			result := struct {
				registerTime time.Duration
				loginTime    time.Duration
				error        error
			}{}

			// Measure registration time
			registerStart := time.Now()
			_, err := client.Register(ctx, &pb.RegisterRequest{
				Username: username,
				Password: "performance_test_pass",
			})
			registerDuration := time.Since(registerStart)

			if err != nil {
				result.error = fmt.Errorf("registration failed: %v", err)
				mu.Lock()
				results[idx] = result
				mu.Unlock()
				return
			}

			result.registerTime = registerDuration

			// Measure login time
			loginStart := time.Now()
			_, err = client.Login(ctx, &pb.LoginRequest{
				Username: username,
				Password: "performance_test_pass",
			})
			loginDuration := time.Since(loginStart)

			if err != nil {
				result.error = fmt.Errorf("login failed: %v", err)
			} else {
				result.loginTime = loginDuration
			}

			mu.Lock()
			results[idx] = result
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	totalDuration := time.Since(startTime)

	// Analyze results
	var totalRegisterTime time.Duration
	var totalLoginTime time.Duration
	var totalCombinedTime time.Duration
	successCount := 0
	errorCount := 0

	for _, result := range results {
		if result.error != nil {
			errorCount++
			t.Logf("Error in performance test: %v", result.error)
		} else {
			successCount++
			totalRegisterTime += result.registerTime
			totalLoginTime += result.loginTime
			totalCombinedTime += result.registerTime + result.loginTime
		}
	}

	// Calculate averages
	if successCount == 0 {
		t.Fatal("All performance test operations failed")
	}

	avgRegisterTime := totalRegisterTime / time.Duration(successCount)
	avgLoginTime := totalLoginTime / time.Duration(successCount)
	avgCombinedTime := totalCombinedTime / time.Duration(successCount)

	t.Logf("Performance Test Results:")
	t.Logf("  - Test duration: %v", totalDuration)
	t.Logf("  - Successful operations: %d/%d", successCount, userCount)
	t.Logf("  - Average registration time: %v", avgRegisterTime)
	t.Logf("  - Average login time: %v", avgLoginTime)
	t.Logf("  - Average combined time: %v", avgCombinedTime)

	// Check against threshold
	if avgCombinedTime > maxAvgResponseTime {
		t.Logf("Warning: Average response time (%v) exceeds threshold (%v)",
			avgCombinedTime, maxAvgResponseTime)
	} else {
		t.Logf("Performance is within acceptable threshold")
	}
}
