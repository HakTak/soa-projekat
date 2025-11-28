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

	// ✅ REŠEN KONFLIKT 1: Koristimo COMMON biblioteke (Incoming grana)
	commonMiddleware "PROJEKAT/COMMON/middleware"
	pbCart "PROJEKAT/COMMON/shopping-cart/proto"
	pbTour "PROJEKAT/COMMON/tour/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	// --- Importi za OpenTelemetry (HEAD grana) ---
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"gorm.io/plugin/opentelemetry/tracing"
)

// --- Inicijalizacija Tracera (Zadržano iz HEAD) ---
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
	// --- 1. POKRETANJE TRACINGA (Zadržano iz HEAD) ---
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
	if err := db.Use(tracing.NewPlugin()); err != nil {
		log.Printf("Failed to use gorm tracing plugin: %v", err)
	}

	repo := repository.NewCartRepository(db)

	// ========================================================================
	// 4. TOUR CLIENT SETUP (SA TRACINGOM)
	// ========================================================================

	// ✅ REŠEN KONFLIKT 2:
	// - Adresa je "tour:8083" (iz Incoming grane, jer je tamo prebačen Tour servis)
	// - Dodajemo 'WithStatsHandler' (iz HEAD grane) da bi Tracing radio
	tourConn, err := grpc.Dial(
		"tour:8083",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // <--- TRACING
	)
	if err != nil {
		log.Fatalf("Failed to connect to Tour Service: %v", err)
	}
	defer tourConn.Close()

	tourClient := pbTour.NewTourServiceClient(tourConn)

	svc := service.NewShoppingCartService(repo, tourClient)
	cartHandler := handlers.NewShoppingCartHandler(svc)

	// ========================================================================
	// 5. START GRPC SERVER
	// ========================================================================

	// ✅ REŠEN KONFLIKT 3: Port 9092 (iz Incoming grane) da se slaže sa Gateway-om
	lis, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Kreiramo server sa Tracingom (Zadržano iz HEAD)
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),

		grpc.UnaryInterceptor(commonMiddleware.MetadataExtractorInterceptor),
	)

	pbCart.RegisterShoppingCartServiceServer(grpcServer, cartHandler)
	reflection.Register(grpcServer)

	log.Println("Shopping Cart Microservice (gRPC) started on port 9092")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
