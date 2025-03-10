package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
	"github.com/golang-jwt/jwt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthHandler manages authentication-related gRPC endpoints
type AuthHandler struct {
	pb.UnimplementedUserManagementServiceServer
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
func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	ctx, span := h.Tracer.Start(ctx, "Register")
	defer span.End()

	// Validate input
	if len(req.Username) < 3 {
		span.RecordError(fmt.Errorf("username too short"))
		span.SetAttributes(attribute.Bool("success", false), attribute.String("error", "username_too_short"))
		return nil, status.Errorf(codes.InvalidArgument, "username must be at least 3 characters")
	}
	if len(req.Password) < 8 {
		span.RecordError(fmt.Errorf("password too short"))
		span.SetAttributes(attribute.Bool("success", false), attribute.String("error", "password_too_short"))
		return nil, status.Errorf(codes.InvalidArgument, "password must be at least 8 characters")
	}

	// Check if username already exists
	var exists bool
	err := h.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", req.Username).Scan(&exists)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false))
		return nil, status.Errorf(codes.Internal, "failed to check username: %v", err)
	}
	if exists {
		span.RecordError(fmt.Errorf("username exists"))
		span.SetAttributes(attribute.Bool("success", false), attribute.String("error", "username_exists"))
		return nil, status.Errorf(codes.AlreadyExists, "username already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false))
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	// Insert the user into the database
	var userID int64
	err = h.DB.QueryRowContext(
		ctx,
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
		req.Username, string(hashedPassword),
	).Scan(&userID)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false))
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	span.SetAttributes(
		attribute.Bool("success", true),
		attribute.Int64("user_id", userID),
	)

	return &pb.RegisterResponse{
		UserId:   userID,
		Username: req.Username,
		Message:  "User registered successfully",
	}, nil
}

// Login handles user login requests
func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	ctx, span := h.Tracer.Start(ctx, "Login")
	defer span.End()

	// Get user from database
	var userID int64
	var username, passwordHash string
	err := h.DB.QueryRowContext(
		ctx,
		"SELECT id, username, password_hash FROM users WHERE username = $1",
		req.Username,
	).Scan(&userID, &username, &passwordHash)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false), attribute.String("error", "user_not_found"))
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	// Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false), attribute.String("error", "invalid_password"))
		return nil, status.Errorf(codes.Unauthenticated, "invalid password")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString(h.JWTSecret)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("success", false))
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	span.SetAttributes(
		attribute.Bool("success", true),
		attribute.Int64("user_id", userID),
	)

	return &pb.LoginResponse{
		Token:  tokenString,
		UserId: userID,
	}, nil
}

// Health handles health check requests
func (h *AuthHandler) Health(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Status: "ok",
	}, nil
}
