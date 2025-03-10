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

	"github.com/MichaelRobotics/Kubernetes/feature/usermanagment-service/opentelemetry-demo/src/db/postgres"

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

// Variables for dependency injection in tests
var (
	resource  *sdkresource.Resource
	sqlOpen   = sql.Open
	osExit    = os.Exit
	logFatalf = log.Fatalf
	logFatal  = log.Fatal
)

// Provider functions to make testing easier
var (
	dbProvider = func() *sql.DB {
		// Use the centralized database module to establish the connection
		dbConn, err := postgres.GetConnectionFromEnv("DB_CONN")
		if err != nil {
			logFatal(err.Error())
			osExit(1)
			return nil // This line is never reached but helps with testing
		}

		// Check if migrations are enabled (disabled by default)
		enableMigrations := false
		if migrationsEnv := os.Getenv("ENABLE_MIGRATIONS"); migrationsEnv == "true" {
			enableMigrations = true
		}

		// Create a user repository
		userRepo := postgres.NewUserRepository(dbConn)

		// Only run migrations if explicitly enabled
		if enableMigrations {
			log.Println("Migrations enabled. Running database migrations...")
			if err := userRepo.EnsureTablesExist(); err != nil {
				logFatalf("Failed to ensure database tables exist: %v", err)
				osExit(1)
				return nil // This line is never reached but helps with testing
			}
			log.Println("Database migrations completed successfully.")
		} else {
			log.Println("Migrations disabled. Skipping database migrations.")
		}

		return dbConn.DB
	}
)

const (
	defaultPort = "8082"
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
	return dbProvider()
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
