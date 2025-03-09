//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate protoc --go_out=./ --go-grpc_out=./ --proto_path=../../pb ../../pb/demo.proto

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "github.com/opentelemetry/demo/src/usermanagementservice/genproto/oteldemo"
)

var (
	resource *sdkresource.Resource
)

const (
	defaultPort = "8080"
)

func initTracer() *sdktrace.TracerProvider {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "otel-collector:4317"
	}

	ctx := context.Background()

	res, err := sdkresource.New(ctx,
		sdkresource.WithAttributes(
			semconv.ServiceNameKey.String("usermanagementservice"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}
	resource = res

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func initDB() *sql.DB {
	dbConn := os.Getenv("DB_CONN")
	if dbConn == "" {
		log.Fatal("DB_CONN environment variable not set")
	}

	var err error
	var db *sql.DB
	db, err = sql.Open("postgres", dbConn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Initialize the database schema
	initStmt := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL
	);
	`
	_, err = db.Exec(initStmt)
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	return db
}

type server struct {
	pb.UnimplementedUserManagementServiceServer
	db        *sql.DB
	tracer    trace.Tracer
	jwtSecret []byte
}

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	ctx, span := s.tracer.Start(ctx, "Register")
	defer span.End()

	// Validate input
	if len(req.Username) < 3 {
		span.RecordError(fmt.Errorf("username too short"))
		return nil, status.Errorf(codes.InvalidArgument, "username must be at least 3 characters")
	}
	if len(req.Password) < 8 {
		span.RecordError(fmt.Errorf("password too short"))
		return nil, status.Errorf(codes.InvalidArgument, "password must be at least 8 characters")
	}

	// Check if username already exists
	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", req.Username).Scan(&exists)
	if err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "failed to check username: %v", err)
	}
	if exists {
		span.RecordError(fmt.Errorf("username already exists"))
		return nil, status.Errorf(codes.AlreadyExists, "username already exists")
	}

	// Hash the password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	// Insert the user into the database
	var userID int32
	err = s.db.QueryRowContext(
		ctx,
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
		req.Username, string(hashedPass),
	).Scan(&userID)
	if err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.RegisterResponse{
		UserId:   userID,
		Username: req.Username,
		Message:  "User registered successfully",
	}, nil
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	ctx, span := s.tracer.Start(ctx, "Login")
	defer span.End()

	// Get user from database
	var userID int32
	var username, passwordHash string
	err := s.db.QueryRowContext(
		ctx,
		"SELECT id, username, password_hash FROM users WHERE username = $1",
		req.Username,
	).Scan(&userID, &username, &passwordHash)
	if err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	// Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.Unauthenticated, "invalid password")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &pb.LoginResponse{
		Token:  tokenString,
		UserId: userID,
	}, nil
}

func (s *server) Health(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Status: "ok",
	}, nil
}

func (s *server) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *server) Watch(req *healthpb.HealthCheckRequest, ws healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}

func main() {
	port := os.Getenv("USER_SVC_URL")
	if port == "" {
		port = defaultPort
	}
	if port[0] == ':' {
		port = port[1:]
	}

	// Initialize tracer
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	tracer := tp.Tracer("usermanagementservice")

	// Initialize database
	db := initDB()
	defer db.Close()

	// Get JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}

	// Create server instance
	srv := &server{
		db:        db,
		tracer:    tracer,
		jwtSecret: []byte(jwtSecret),
	}

	// Set up gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	// Register services
	pb.RegisterUserManagementServiceServer(grpcServer, srv)
	healthpb.RegisterHealthServer(grpcServer, srv)
	reflection.Register(grpcServer)

	// Start server
	log.Printf("user management service listening on port %s", port)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down user management service...")
	grpcServer.GracefulStop()
} 