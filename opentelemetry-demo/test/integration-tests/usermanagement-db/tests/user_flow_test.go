package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client pb.UserManagementServiceClient

func init() {
	// Establish gRPC connection
	conn, err := grpc.Dial("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to gRPC service: %v", err))
	}
	// Note: We're not closing the connection as tests will use it
	client = pb.NewUserManagementServiceClient(conn)
}

func TestRegisterLoginFlow(t *testing.T) {
	ctx := context.Background()

	// Generate a unique username to avoid conflicts
	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	password := "testpassword"

	// Step 1: Register a new user
	registerReq := &pb.RegisterRequest{
		Username: username,
		Password: password,
	}
	registerResp, err := client.Register(ctx, registerReq)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if registerResp.GetUserId() == 0 {
		t.Error("Expected non-zero user ID in RegisterResponse")
	}
	if registerResp.GetUsername() != username {
		t.Errorf("Expected username %q, got %q", username, registerResp.GetUsername())
	}
	if registerResp.GetMessage() == "" {
		t.Error("Expected a non-empty message in RegisterResponse")
	}

	// Step 2: Login with the same credentials
	loginReq := &pb.LoginRequest{
		Username: username,
		Password: password,
	}
	loginResp, err := client.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if loginResp.GetToken() == "" {
		t.Error("Expected a non-empty token in LoginResponse")
	}
	if loginResp.GetUserId() != registerResp.GetUserId() {
		t.Errorf("Expected user ID %d from Login to match Register, got %d", registerResp.GetUserId(), loginResp.GetUserId())
	}

	// Step 3: "Retrieve User" (assumed implicit in LoginResponse)
	// Since no GetUser method is provided, verify user_id consistency
	t.Logf("User registered and logged in successfully with user_id: %d and token: %s", loginResp.GetUserId(), loginResp.GetToken())
}
