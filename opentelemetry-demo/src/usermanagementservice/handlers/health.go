package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthHandler handles health check requests
type HealthHandler struct {}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health handles health check requests
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Note: The Health gRPC handler is implemented in auth.go as a part of the AuthHandler
// This allows us to maintain a single server implementation while still having modular code 