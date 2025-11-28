package main

import (
	"context"
	"log"
	"net"
	"os"

	"tour-service/internal/api"
	"tour-service/internal/model"
	"tour-service/internal/repository"
	"tour-service/internal/service"

	// Koristimo Common pakete iz incoming grane
	commonMiddleware "PROJEKAT/COMMON/middleware"
	pbShop "PROJEKAT/COMMON/shopping-cart/proto"
	pb "PROJEKAT/COMMON/tour/proto"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// --- Importi za OpenTelemetry (Zadržavamo iz HEAD) ---
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"gorm.io/plugin/opentelemetry/tracing" // Plugin za Gorm
)

// --- Funkcija za inicijalizaciju Jaegera (Zadržavamo iz HEAD) ---
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
		serviceName = "tour-service"
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
	// --- 1. POKRETANJE TRACINGA (Zadržavamo iz HEAD) ---
	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}

	// --- 2. DODAVANJE TRACINGA ZA BAZU (Zadržavamo iz HEAD) ---
	if err := db.Use(tracing.NewPlugin()); err != nil {
		log.Printf("Failed to use gorm tracing plugin: %v", err)
	}

	// Auto migrate tables
	db.AutoMigrate(&model.Tour{}, &model.Keypoint{}, &model.Review{}, &model.RouteOption{})

	// Repository -> Service -> Handler
	repo := repository.NewTourRepository(db)
	reviewRepo := repository.CreateReviewRepository(db)
	repoTourExecution := repository.NewTourExecutionRepository(db)
	svc := service.NewTourService(repo, reviewRepo, repoTourExecution)

	// --- REŠENJE KONFLIKTA KOD SERVERA ---
	// Koristimo strukturu iz Incoming grane (NewTourGRPCServer umesto HTTP handlera),
	// ali dodajemo Tracing iz HEAD grane.

	shopConn, err := grpc.NewClient(
		"shopping-cart:9092",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)

	if err != nil {
		log.Fatalf("Failed to connect to stakeholders: %v", err)
	}
	defer shopConn.Close()

	shopClient := pbShop.NewShoppingCartServiceClient(shopConn)

	handler := api.NewTourGRPCServer(svc, shopClient)

	// Incoming grana koristi port 8083 za gRPC
	lis, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// --- SPAJANJE INTERCEPTORA ---
	grpcServer := grpc.NewServer(
		// 1. Dodajemo Tracing (iz HEAD)
		grpc.StatsHandler(otelgrpc.NewServerHandler()),

		// 2. Dodajemo Metadata Middleware (iz Incoming)
		grpc.UnaryInterceptor(commonMiddleware.MetadataExtractorInterceptor),
	)

	pb.RegisterTourServiceServer(grpcServer, handler)

	log.Println("Tour gRPC service running on port 8083 (with Tracing)")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
