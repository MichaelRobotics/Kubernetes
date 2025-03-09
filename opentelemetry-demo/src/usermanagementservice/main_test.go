package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func TestHealthChecker_Check(t *testing.T) {
	healthChecker := &HealthChecker{}

	resp, err := healthChecker.Check(context.Background(), &healthpb.HealthCheckRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.Status)
}

func TestHealthChecker_Watch(t *testing.T) {
	healthChecker := &HealthChecker{}

	err := healthChecker.Watch(&healthpb.HealthCheckRequest{}, nil)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok, "expected a gRPC status error")
	assert.Equal(t, codes.Unimplemented, st.Code(), "expected Unimplemented error code")
}
