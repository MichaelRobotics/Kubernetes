package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/opentelemetry/demo/src/usermanagementservice/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler manages authentication-related HTTP endpoints
type AuthHandler struct {
	DB        *sql.DB
	Tracer    trace.Tracer
	JWTSecret []byte
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(db *sql.DB, tracer trace.Tracer, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{
		DB:        db,
		Tracer:    tracer,
		JWTSecret: jwtSecret,
	}
}

// Register handles user registration requests
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.Tracer.Start(r.Context(), "RegisterUser")
	defer span.End()

	// Parse request body
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false))
		return
	}

	// Validate input
	if len(creds.Username) < 3 {
		respondWithError(w, "Username must be at least 3 characters", http.StatusBadRequest)
		span.SetAttributes(attribute.Bool("success", false), attribute.String("error", "username_too_short"))
		return
	}
	if len(creds.Password) < 8 {
		respondWithError(w, "Password must be at least 8 characters", http.StatusBadRequest)
		span.SetAttributes(attribute.Bool("success", false), attribute.String("error", "password_too_short"))
		return
	}

	// Check if username already exists
	var exists bool
	err := h.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", creds.Username).Scan(&exists)
	if err != nil {
		log.Printf("Database error: %v", err)
		respondWithError(w, "Internal server error", http.StatusInternalServerError)
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false))
		return
	}
	if exists {
		respondWithError(w, "Username already exists", http.StatusConflict)
		span.SetAttributes(attribute.Bool("success", false), attribute.String("error", "username_exists"))
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		respondWithError(w, "Internal server error", http.StatusInternalServerError)
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false))
		return
	}

	// Insert user into database
	var userID int
	err = h.DB.QueryRowContext(
		ctx,
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
		creds.Username, string(hashedPassword),
	).Scan(&userID)
	if err != nil {
		log.Printf("User creation error: %v", err)
		respondWithError(w, "Failed to create user", http.StatusInternalServerError)
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false))
		return
	}

	span.SetAttributes(
		attribute.Bool("success", true),
		attribute.Int("user_id", userID),
	)

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":  userID,
		"username": creds.Username,
		"message":  "User registered successfully",
	})
}

// Login handles user login requests
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.Tracer.Start(r.Context(), "LoginUser")
	defer span.End()

	// Parse request body
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, "Invalid request payload", http.StatusBadRequest)
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false))
		return
	}

	// Retrieve user from database
	var user models.User
	var passwordHash string
	err := h.DB.QueryRowContext(
		ctx,
		"SELECT id, username, password_hash FROM users WHERE username = $1",
		creds.Username,
	).Scan(&user.ID, &user.Username, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, "Invalid credentials", http.StatusUnauthorized)
			span.SetAttributes(
				attribute.Bool("success", false),
				attribute.String("error", "invalid_credentials"),
			)
			return
		}
		log.Printf("Database error: %v", err)
		respondWithError(w, "Internal server error", http.StatusInternalServerError)
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false))
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(creds.Password)); err != nil {
		respondWithError(w, "Invalid credentials", http.StatusUnauthorized)
		span.SetAttributes(
			attribute.Bool("success", false),
			attribute.String("error", "invalid_credentials"),
		)
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString(h.JWTSecret)
	if err != nil {
		log.Printf("Token signing error: %v", err)
		respondWithError(w, "Internal server error", http.StatusInternalServerError)
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false))
		return
	}

	span.SetAttributes(
		attribute.Bool("success", true),
		attribute.Int("user_id", user.ID),
	)

	// Return token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.LoginResponse{
		Token:  tokenString,
		UserID: user.ID,
	})
}

// Helper function to send error responses
func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
} 