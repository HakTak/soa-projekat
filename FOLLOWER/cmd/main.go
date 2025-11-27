package main

import (
	"context"
	"log"
	"net"
	"os"

	"FOLLOWER/internal/config"
	"FOLLOWER/internal/database"
	"FOLLOWER/internal/handlers"
	"FOLLOWER/internal/repository"
	"FOLLOWER/internal/service"

	pb "PROJEKAT/COMMON/follower/proto"
	commonMiddleware "PROJEKAT/COMMON/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// --- Importi za OpenTelemetry ---
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// --- Funkcija za inicijalizaciju Tracera ---
func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// Cita adresu Jaegera (iz docker-compose)
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "jaeger:4317"
	}

	// Kreira exporter ka Jaegeru (gRPC)
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("jaeger:4317"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "follower-service"
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

func main() {
	// --- 1. POKRETANJE TRACINGA ---
	tp, err := initTracer()
	if err != nil {
		log.Fatalf("Failed to init tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	cfg := config.LoadConfig()

	driver, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Cannot connect to Neo4j: %v", err)
	}
	defer driver.Close(context.Background())

	repo := repository.NewfollowRepository(driver)

	log.Println("Applying Neo4j constraints...")
	if err := repo.EnsureConstraints(context.Background()); err != nil {
		log.Printf("Warning: Could not ensure constraints: %v", err)
	} else {
		log.Println("Constraints applied.")
	}

	svc := service.NewFollowService(repo)

	followerHandler := handlers.NewFollowerHandler(svc)

	//start grpc server port 9090
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// --- 2. SETUP GRPC SERVERA ---
	grpcServer := grpc.NewServer(
		// A) Tracing Handler (Ovo hvata requeste za Jaeger)
		grpc.StatsHandler(otelgrpc.NewServerHandler()),

		// B) Tvoj postojeci Middleware (Ovo vadi metadata)
		grpc.UnaryInterceptor(commonMiddleware.MetadataExtractorInterceptor),
	)

	pb.RegisterFollowerServiceServer(grpcServer, followerHandler)

	reflection.Register(grpcServer)

	log.Println("Follower Microservice (gRPC) started on port 9090")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
