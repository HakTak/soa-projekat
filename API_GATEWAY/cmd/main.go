package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"PROJEKAT/API_GATEWAY/middleware"
	pbAuth "PROJEKAT/COMMON/auth/proto"
	pbBlog "PROJEKAT/COMMON/blog/proto"
	pbFollower "PROJEKAT/COMMON/follower/proto"
	pbShoppingCart "PROJEKAT/COMMON/shopping-cart/proto"
	pbStakeholders "PROJEKAT/COMMON/stakeholders/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	// --- DODATI IMPORTI ZA TRACING ---
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// --- OVO JE JEDINA NOVA FUNKCIJA (BOILERPLATE) ---
// Ona sluzi samo da se povezemo na Jaeger kontejner.
func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// Povezivanje na Jaeger
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("jaeger:4317"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("api-gateway"), // Ime servisa
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
	// Umesto middleware.InitTracer, koristimo ovu lokalnu funkciju
	// da budemo sigurni da gadja pravi Jaeger port.
	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 1. GRPC GATEWAY MUX
	gwmux := runtime.NewServeMux()

	// --- OVO JE KLJUCNA IZMENA ---
	// Dodajemo 'grpc.WithStatsHandler(otelgrpc.NewClientHandler())'
	// Ovo omogucava da Gateway prosledi TraceID ka mikroservisima.
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // <--- DODATO
	}

	// Registracija AUTH servisa
	err = pbAuth.RegisterAuthServiceHandlerFromEndpoint(ctx, gwmux, "auth:8082", opts)
	if err != nil {
		log.Fatalf("Failed to register Auth: %v", err)
	}

	// Registracija STAKEHOLDERS servisa
	err = pbStakeholders.RegisterStakeholdersServiceHandlerFromEndpoint(ctx, gwmux, "stakeholders:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register Stakeholders: %v", err)
	}

	// Registracija FOLLOWER servisa
	err = pbFollower.RegisterFollowerServiceHandlerFromEndpoint(ctx, gwmux, "follower:9090", opts)
	if err != nil {
		log.Fatalf("Faild to register Follower: %v", err)
	}

	// Registracija Blog servisa
	err = pbBlog.RegisterBlogServiceHandlerFromEndpoint(ctx, gwmux, "blog:9091", opts)
	if err != nil {
		log.Fatalf("Failed to register Blog: %v", err)
	}

	// Registracija Shopping Cart servisa
	err = pbShoppingCart.RegisterShoppingCartServiceHandlerFromEndpoint(ctx, gwmux, "shopping-cart:9092", opts)
	if err != nil {
		log.Fatalf("Failed to register Shopping cart: %v", err)
	}

	// 2. GLAVNI RUTER (Standardni HTTP)
	rootMux := http.NewServeMux()

	// Sve rute saljemo na gRPC Gateway
	rootMux.Handle("/", gwmux)

	// 3. AUTH MIDDLEWARE
	authHandler := middleware.AuthMiddleware(rootMux)

	// 4. CORS MIDDLEWARE
	finalHandler := middleware.CorsMiddleware(authHandler)

	// 5. TRACING MIDDLEWARE
	// Ovo hvata dolazni zahtev od Frontenda
	tracingHandler := otelhttp.NewHandler(finalHandler, "api-gateway-request")

	fmt.Println("API Gateway running on port 8080...")
	// Pokrecemo tracingHandler umesto finalHandler
	if err := http.ListenAndServe(":8080", tracingHandler); err != nil {
		log.Fatal(err)
	}
}
