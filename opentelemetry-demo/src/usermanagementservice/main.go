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

	"github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/handlers"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
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

// HealthChecker implements the gRPC health check service
type HealthChecker struct {
	healthpb.UnimplementedHealthServer
}

func (s *HealthChecker) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *HealthChecker) Watch(req *healthpb.HealthCheckRequest, ws healthpb.Health_WatchServer) error {
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

	// Create handlers
	authHandler := handlers.NewAuthHandler(db, tracer, []byte(jwtSecret))
	healthChecker := &HealthChecker{}

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
	pb.RegisterUserManagementServiceServer(grpcServer, authHandler)
	healthpb.RegisterHealthServer(grpcServer, healthChecker)
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
