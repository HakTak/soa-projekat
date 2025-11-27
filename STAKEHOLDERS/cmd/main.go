package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"stakeholders/internal/handlers"
	"stakeholders/internal/model"
	"stakeholders/internal/repository"
	"stakeholders/internal/service"

	commonMiddleware "PROJEKAT/COMMON/middleware" // Tvoj middleware
	pb "PROJEKAT/COMMON/stakeholders/proto"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// --- Importi za OpenTelemetry (Jaeger) ---
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"gorm.io/plugin/opentelemetry/tracing"
)

// --- Funkcija za inicijalizaciju Jaegera ---
func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// Cita adresu iz docker-compose ili koristi default
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "jaeger:4317"
	}

	// Kreiramo exporter (gRPC komunikacija sa Jaegerom)
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("jaeger:4317"),
		otlptracegrpc.WithInsecure(), // Bitno za Docker (bez SSL-a)
	)
	if err != nil {
		return nil, err
	}

	// Ime servisa koje ce pisati u Grafani
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "stakeholders-service"
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

	// Postavljamo globalne propagatore (da prenosimo TraceID kroz servise)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

func main() {
	// --- 1. POKRETANJE TRACINGA ---
	tp, err := initTracer()
	if err != nil {
		log.Fatalf("failed to init tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=stakeholders password=secret dbname=stakeholders port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	// --- 2. DODAVANJE TRACINGA ZA BAZU ---
	if err := db.Use(tracing.NewPlugin()); err != nil {
		log.Printf("Failed to use gorm tracing plugin: %v", err)
	}

	db.AutoMigrate(&model.Profile{})
	// seedProfiles(db)

	repo := repository.NewGormProfileRepo(db)
	svc := service.NewProfileService(repo)

	profileHandler := handlers.NewProfileHandler(svc)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	// --- 3. KONFIGURACIJA GRPC SERVERA SA TRACINGOM ---
	// Koristimo ChainUnaryInterceptor da povezemo OpenTelemetry I tvoj Metadata extractor
	grpcServer := grpc.NewServer(
		// 1. OpenTelemetry Tracing (hvata sve zahteve automatski)
		grpc.StatsHandler(otelgrpc.NewServerHandler()),

		// 2. Tvoj Middleware (ostaje kao Interceptor)
		grpc.UnaryInterceptor(commonMiddleware.MetadataExtractorInterceptor),
	)

	pb.RegisterStakeholdersServiceServer(grpcServer, profileHandler)

	fmt.Println("Stakeholders gRPC Service running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

/*
package main

import (
	"log"
	"net"
	stdhttp "net/http"
	"os"

	"stakeholders/internal/api/http"
	"stakeholders/internal/model"
	"stakeholders/internal/repository"
	"stakeholders/internal/service"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	pb "stakeholders/common/genproto"
	grpcserver "stakeholders/internal/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=stakeholders password=secret dbname=stakeholders port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// Automigrate profile table (creates table if missing). Safe for dev; consider controlled migrations in prod.
	if err := db.AutoMigrate(&model.Profile{}); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}

	var count int64
	db.Model(&model.Profile{}).Count(&count)
	if count == 0 {
		seedProfiles(db)
	}

	repo := repository.NewGormProfileRepo(db)
	svc := service.NewProfileService(repo)
	h := http.NewHandler(svc)

	go func() {
		grpcAddr := ":50051"
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("failed to listen grpc: %v", err)
		}

		grpcSrv := grpc.NewServer()
		pb.RegisterStakeholdersServer(grpcSrv, grpcserver.NewServer(svc))

		reflection.Register(grpcSrv)

		log.Printf("gRPC listening on %s", grpcAddr)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	r := chi.NewRouter()
	r.Use(http.AuthMiddleware)
	h.RegisterRoutes(r)

	addr := ":8081"
	log.Printf("listening on %s", addr)
	if err := stdhttp.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func seedProfiles(db *gorm.DB) {
	profiles := []model.Profile{
		{
			UserID:         "d6daf005-891f-4ae2-8a89-55ed4770a532",
			FirstName:      "Alice",
			LastName:       "Smith",
			ProfilePicture: "",
			Biography:      "Tour guide from NYC.",
			Motto:          "Adventure awaits!",
			Role:           model.RoleGuide,
			IsBlocked:      false,
		},
		{
			UserID:         "831c17e1-9be9-4bc8-a2d5-555a36ceac46",
			FirstName:      "Bob",
			LastName:       "Brown",
			ProfilePicture: "",
			Biography:      "Tourist from CA.",
			Motto:          "Live, laugh, travel.",
			Role:           model.RoleTourist,
			IsBlocked:      false,
		},
		{
			UserID:         "123e4567-e89b-12d3-a456-426614174003",
			FirstName:      "Carol",
			LastName:       "White",
			ProfilePicture: "",
			Biography:      "Admin user.",
			Motto:          "Keeping things running.",
			Role:           model.RoleAdmin,
			IsBlocked:      false,
		},
	}

	result := db.Create(&profiles)
	if result.Error != nil {
		log.Fatalf("failed to seed profiles: %v", result.Error)
	}

	log.Println("Seeded profiles successfully")
}
*/
