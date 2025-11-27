package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"tour-service/internal/api"
	"tour-service/internal/model"
	"tour-service/internal/repository"
	"tour-service/internal/service"
	pb "tour-service/protobuf"

	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi" // Middleware za Chi
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing" // Plugin za Gorm

	"google.golang.org/grpc"

	// --- Importi za OpenTelemetry ---
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// --- Funkcija za inicijalizaciju Jaegera (OpenTelemetry) ---
func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// Cita adresu Jaegera iz environment varijable ili koristi default
	// U docker-compose smo stavili OTEL_EXPORTER_OTLP_ENDPOINT=http://jaeger:4317
	// Go exporter ochekuje host:port format za gRPC
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "jaeger:4317"
	}

	// Kreiramo exporter koji salje podatke Jaegeru preko gRPC-a
	// Insecure je potreban jer unutar Dockera ne koristimo HTTPS sertifikate
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("jaeger:4317"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// Definisemo ime servisa (ovo ce pisati u Jaeger UI)
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

	// Kreiramo Tracer Provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Postavljamo globalni tracer i propagator (bitno za povezivanje servisa)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

func main() {
	// --- 1. POKRETANJE TRACINGA ---
	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	// Osiguravamo da se podaci posalju pre gasenja servisa
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

	// --- 2. DODAVANJE TRACINGA ZA BAZU (GORM) ---
	if err := db.Use(tracing.NewPlugin()); err != nil {
		log.Printf("Failed to use gorm tracing plugin: %v", err)
	}

	// Auto migrate tables
	db.AutoMigrate(&model.Tour{}, &model.Keypoint{}, &model.Review{}, &model.RouteOption{})

	// Repository -> Service -> Handler
	repo := repository.NewTourRepository(db)
	reviewRepo := repository.CreateReviewRepository(db)
	svc := service.NewTourService(repo, reviewRepo)
	handler := api.NewTourHandler(svc)

	go func() {
		lis, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatalf("failed to listen for gRPC: %v", err)
		}

		// --- 3. DODAVANJE TRACINGA ZA gRPC SERVER ---
		// StatsHandler automatski prati sve gRPC pozive
		grpcServer := grpc.NewServer(
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
		)

		grpcHandler := api.NewTourGrpcHandler(svc)
		pb.RegisterTourServiceServer(grpcServer, grpcHandler)

		log.Println("gRPC Server started on port :50052")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	r := chi.NewRouter()

	// --- 4. DODAVANJE TRACINGA ZA HTTP (Chi Router) ---
	// "tour-service" je ime koje ce se videti u spanovima
	r.Use(otelchi.Middleware("tour-service"))

	// CORS MIDDLEWARE
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	r.Post("/tour", handler.CreateTour)
	r.Get("/tour/{id}", handler.GetTour)
	r.Get("/tours", handler.GetAllTours)
	r.Get("/tours/myTours/{userId}", handler.GetToursByUser)
	r.Delete("/tour/{id}", handler.DeleteTour)
	r.Patch("/tour/update", handler.UpdateTour)

	r.Post("/review/create", handler.CreateReview)
	r.Get("/review/tour/{tourId}", handler.GetReviewsByTour)
	r.Get("/review/getAll", handler.GetAllReviews)
	r.Delete("/review/delete/{id}", handler.DeleteReview)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Println("Tour service running on port", port)
	http.ListenAndServe(":"+port, r)
}
