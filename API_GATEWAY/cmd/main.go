package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"PROJEKAT/API_GATEWAY/middleware"
	pbAuth "PROJEKAT/COMMON/auth/proto"
	pbBlog "PROJEKAT/COMMON/blog/proto"
	pbFollower "PROJEKAT/COMMON/follower/proto"
	pbShoppingCart "PROJEKAT/COMMON/shopping-cart/proto"
	pbStakeholders "PROJEKAT/COMMON/stakeholders/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp" // <--- NOVA BIBLIOTEKA
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// --- 0. TRACING SETUP (DODATO) ---
	// Postavljamo adresu Jaegera (HTTP Collector port)
	os.Setenv("JAEGER_ENDPOINT", "http://jaeger:14268/api/traces")

	// Inicijalizujemo Tracer iz tvog middleware paketa
	tp, err := middleware.InitTracer()
	if err != nil {
		log.Fatal(err)
	}

	// Osiguravamo da se podaci pošalju pre gašenja servisa (dobra praksa)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Postavljamo globalni tracer provider
	otel.SetTracerProvider(tp)
	// Postavljamo propagator da bi se Trace ID prenosio kroz servise
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	// ---------------------------------

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 1. GRPC GATEWAY MUX
	gwmux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

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

	// 5. TRACING MIDDLEWARE (DODATO NA KRAJU)
	// Obmotavamo ceo handler sa OpenTelemetry instrumentacijom.
	// "api-gateway-request" je ime operacije koje ćeš videti u Jaegeru.
	tracingHandler := otelhttp.NewHandler(finalHandler, "api-gateway-request")

	fmt.Println("API Gateway running on port 8080...")
	// Ovde sada pokrećemo tracingHandler umesto finalHandler
	if err := http.ListenAndServe(":8080", tracingHandler); err != nil {
		log.Fatal(err)
	}
}
