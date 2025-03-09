package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// User represents the user model stored in the database
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // Never expose in JSON
}

// Credentials represents the login/registration request body
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response with JWT token
type LoginResponse struct {
	Token string `json:"token"`
	UserID int   `json:"user_id"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	db         *sql.DB
	tracer     trace.Tracer
	jwtSecret  []byte
	serviceURL string
)

func initTracer() *sdktrace.TracerProvider {
	ctx := context.Background()

	// OTLP exporter
	conn, err := grpc.DialContext(ctx,
		os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("failed to create gRPC connection: %v", err)
	}

	// Create OTLP exporter
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatalf("failed to create trace exporter: %v", err)
	}

	// Create resource with service information
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("usermanagementservice"),
		semconv.ServiceVersion("1.0.0"),
	)

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp
}

func initDB() {
	var err error
	connStr := os.Getenv("DB_CONN")
	if connStr == "" {
		connStr = "postgresql://postgres:postgres@postgres:5432/users?sslmode=disable"
	}

	// Connect to PostgreSQL
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Create users table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("failed to create users table: %v", err)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "RegisterUser")
	defer span.End()

	// Parse request body
	var creds Credentials
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
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", creds.Username).Scan(&exists)
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
	err = db.QueryRowContext(
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "LoginUser")
	defer span.End()

	// Parse request body
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, "Invalid request payload", http.StatusBadRequest)
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false))
		return
	}

	// Retrieve user from database
	var user User
	var passwordHash string
	err := db.QueryRowContext(
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

	tokenString, err := token.SignedString(jwtSecret)
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
	json.NewEncoder(w).Encode(LoginResponse{
		Token:  tokenString,
		UserID: user.ID,
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func main() {
	// Initialize tracer
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	tracer = tp.Tracer("usermanagementservice")

	// Initialize database
	initDB()
	defer db.Close()

	// Get JWT secret from environment
	jwtSecretEnv := os.Getenv("JWT_SECRET")
	if jwtSecretEnv == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	jwtSecret = []byte(jwtSecretEnv)

	// Get service URL
	serviceURL = os.Getenv("USER_SVC_URL")
	if serviceURL == "" {
		serviceURL = ":8080"
	}

	// Create router with OpenTelemetry instrumentation
	r := mux.NewRouter()
	r.Use(otelmux.Middleware("usermanagementservice"))

	// Define routes
	r.HandleFunc("/register", registerHandler).Methods("POST")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/health", healthHandler).Methods("GET")

	// Start HTTP server
	http.Handle("/", r)
	wrappedHandler := otelhttp.NewHandler(r, "http.server")

	log.Printf("User Management Service is running on %s", serviceURL)
	log.Fatal(http.ListenAndServe(serviceURL, wrappedHandler))
} 