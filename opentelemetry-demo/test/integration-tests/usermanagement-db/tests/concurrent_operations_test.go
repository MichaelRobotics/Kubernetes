package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
)

// TestConcurrentOperations verifies that the system can handle multiple concurrent user operations
func TestConcurrentOperations(t *testing.T) {
	ctx := context.Background()
	var wg sync.WaitGroup
	concurrentUsers := 5 // Adjust based on your system capacity
	errorChan := make(chan error, concurrentUsers)

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			username := fmt.Sprintf("concurrent_user_%d_%d", time.Now().UnixNano(), idx)

			// Register
			_, err := client.Register(ctx, &pb.RegisterRequest{
				Username: username,
				Password: "concurrent_pass",
			})
			if err != nil {
				errorChan <- fmt.Errorf("concurrent registration %d failed: %v", idx, err)
				return
			}

			// Login
			loginResp, err := client.Login(ctx, &pb.LoginRequest{
				Username: username,
				Password: "concurrent_pass",
			})
			if err != nil {
				errorChan <- fmt.Errorf("concurrent login %d failed: %v", idx, err)
				return
			}

			// Verify token is not empty
			if loginResp.GetToken() == "" {
				errorChan <- fmt.Errorf("concurrent login %d returned empty token", idx)
			}
		}(i)
	}

	wg.Wait()
	close(errorChan)

	// Check if there were any errors
	errCount := 0
	for err := range errorChan {
		t.Error(err)
		errCount++
	}

	if errCount == 0 {
		t.Logf("Successfully processed %d concurrent users", concurrentUsers)
	}
}
