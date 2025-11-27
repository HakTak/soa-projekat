package main

import (
	"context"
	"log"
	"net"
	"os"

	"SHOPPING-CART/internal/database"
	"SHOPPING-CART/internal/handlers"
	"SHOPPING-CART/internal/repository"
	"SHOPPING-CART/internal/service"

	pbCart "PROJEKAT/COMMON/shopping-cart/proto"
	pbTour "SHOPPING-CART/common/genproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

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

// --- Inicijalizacija Tracera (Boilerplate kod) ---
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
		serviceName = "shopping-cart-service"
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

	// 2. Database Connection
	db := database.Connect()

	// --- 3. BAZA TRACING (GORM) ---
	// Dodajemo plugin da vidimo SQL upite u Jaegeru
	if err := db.Use(tracing.NewPlugin()); err != nil {
		log.Printf("Failed to use gorm tracing plugin: %v", err)
	}

	repo := repository.NewCartRepository(db)

	// ========================================================================
	// 4. TOUR CLIENT SETUP (SA TRACINGOM)
	// Ovo je kljucno: Kada Cart zove Tour, moramo poslati TraceID dalje.
	// ========================================================================

	tourConn, err := grpc.Dial(
		"tour:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),

		// --- ISPRAVKA ---
		// Koristimo "WithStatsHandler" jer je ovo klijent (Dial), a ne server.
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to Tour Service: %v", err)
	}
	defer tourConn.Close()

	tourClient := pbTour.NewTourServiceClient(tourConn)

	svc := service.NewShoppingCartService(repo, tourClient)
	cartHandler := handlers.NewShoppingCartHandler(svc)

	// ========================================================================
	// 5. START GRPC SERVER (SA TRACINGOM)
	// ========================================================================

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Kreiramo server sa StatsHandler-om koji automatski prati sve dolazne zahteve
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()), // <--- KLJUCNO: Tracing dolaznih zahteva
	)

	pbCart.RegisterShoppingCartServiceServer(grpcServer, cartHandler)
	reflection.Register(grpcServer)

	log.Println("Shopping Cart Microservice (gRPC) started on port 50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
