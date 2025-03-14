package main

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client pb.UserManagementServiceClient

func TestMain(m *testing.M) {
	// Setup before any tests run
	log.Println("Setting up integration test environment")

	// Establish gRPC connection
	conn, err := grpc.Dial("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC service: %v", err)
	}
	defer conn.Close()

	client = pb.NewUserManagementServiceClient(conn)

	// Wait for the service to be healthy
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for {
		resp, err := client.Health(ctx, &pb.HealthRequest{})
		if err == nil && resp.GetStatus() == "ok" {
			log.Println("Service is healthy, proceeding with tests")
			break
		}
		select {
		case <-ctx.Done():
			log.Fatalf("Service did not become healthy within 30 seconds: %v", err)
		default:
			time.Sleep(1 * time.Second)
		}
	}

	// Run the tests
	exitCode := m.Run()

	// Cleanup after all tests have run
	log.Println("Cleaning up integration test environment")

	os.Exit(exitCode)
}

// TestHealthEndpoint verifies that the health endpoint is working properly
func TestHealthEndpoint(t *testing.T) {
	ctx := context.Background()
	resp, err := client.Health(ctx, &pb.HealthRequest{})
	if err != nil {
		t.Fatalf("Failed to call Health: %v", err)
	}
	if resp.GetStatus() != "ok" {
		t.Errorf("Expected status 'ok', got %q", resp.GetStatus())
	}
}
