package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"auth/internal/database"
	"auth/internal/handlers"
	"auth/internal/repositories"
	"auth/internal/services"

	pbAuth "PROJEKAT/COMMON/auth/proto"
	commonMiddleware "PROJEKAT/COMMON/middleware"
	pbStakeholders "PROJEKAT/COMMON/stakeholders/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// --- Importi za OpenTelemetry ---
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"gorm.io/plugin/opentelemetry/tracing"
)

// --- Inicijalizacija Tracera ---
func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "jaeger:4317"
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("jaeger:4317"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "auth-service"
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

	// Konekcija na bazu
	db := database.Connect()

	// --- 2. BAZA TRACING ---
	if err := db.Use(tracing.NewPlugin()); err != nil {
		log.Printf("Failed to use gorm tracing plugin: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	jwtService := services.NewJWTService("my_super_duper_secret_mega_gg_key_123", "AUTH_SERVICE")

	// --- 3. KLIJENT KA STAKEHOLDERS (Sa Tracingom) ---
	// Kada Auth zove Stakeholders, zelimo da se trace nastavi
	stakeholdersConn, err := grpc.NewClient(
		"stakeholders:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// Dodajemo ClientHandler
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to stakeholders: %v", err)
	}
	defer stakeholdersConn.Close()

	stakeholdersClient := pbStakeholders.NewStakeholdersServiceClient(stakeholdersConn)

	userHandler := handlers.NewUserHandler(userService, jwtService, stakeholdersClient)

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// --- 4. SERVER SETUP (Sa Tracingom) ---
	grpcServer := grpc.NewServer(
		// A) Tracing Handler (za dolazne zahteve)
		grpc.StatsHandler(otelgrpc.NewServerHandler()),

		// B) Tvoj Metadata Interceptor
		grpc.UnaryInterceptor(commonMiddleware.MetadataExtractorInterceptor),
	)

	pbAuth.RegisterAuthServiceServer(grpcServer, userHandler)

	fmt.Println("Auth gRPC Service running on port 8082...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
