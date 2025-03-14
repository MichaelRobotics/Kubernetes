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
	"google.golang.org/grpc/reflection"

	"github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/db/postgres"
	pb "github.com/MichaelRobotics/Kubernetes/opentelemetry-demo/src/usermanagementservice/genproto/oteldemo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
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

		// Migration code has been removed

		return dbConn.DB
	}
)

const (
	defaultPort = "8082"
)

// HealthChecker implements the gRPC health check service
type HealthChecker struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (s *HealthChecker) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (s *HealthChecker) Watch(req *grpc_health_v1.HealthCheckRequest, ws grpc_health_v1.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}

// UserManagementServiceServer combines the AuthHandler with the Health method
type UserManagementServiceServer struct {
	authHandler *handlers.AuthHandler
	pb.UnimplementedUserManagementServiceServer
}

// Register forwards the register request to the auth handler
func (s *UserManagementServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return s.authHandler.Register(ctx, req)
}

// Login forwards the login request to the auth handler
func (s *UserManagementServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return s.authHandler.Login(ctx, req)
}

// Health handles the Health RPC call for the UserManagementService
func (s *UserManagementServiceServer) Health(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	// Use the same status as the gRPC health check
	return &pb.HealthResponse{
		Status: "ok",
	}, nil
}

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

	// Get tracer
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

	// Create the combined service server
	userManagementServer := &UserManagementServiceServer{
		authHandler: authHandler,
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
	pb.RegisterUserManagementServiceServer(grpcServer, userManagementServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthChecker)
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
