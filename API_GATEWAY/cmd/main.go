package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"PROJEKAT/API_GATEWAY/middleware"
	pbAuth "PROJEKAT/COMMON/auth/proto"
	pbStakeholders "PROJEKAT/COMMON/stakeholders/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 1. GRPC GATEWAY MUX
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Registracija AUTH servisa
	err := pbAuth.RegisterAuthServiceHandlerFromEndpoint(ctx, gwmux, "auth:8082", opts)
	if err != nil {
		log.Fatalf("Failed to register Auth: %v", err)
	}

	// Registracija STAKEHOLDERS servisa
	err = pbStakeholders.RegisterStakeholdersServiceHandlerFromEndpoint(ctx, gwmux, "stakeholders:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register Stakeholders: %v", err)
	}

	// 2. GLAVNI RUTER (Standardni HTTP)
	// Ovo koristimo da bi mogli lako da dodamo Swagger ili Health checkove u buducnosti
	rootMux := http.NewServeMux()

	// Sve rute saljemo na gRPC Gateway
	rootMux.Handle("/", gwmux)

	// 3. AUTH MIDDLEWARE
	authHandler := middleware.AuthMiddleware(rootMux)

	//4. CORS MIDDLEWARE
	finalHandler := middleware.CorsMiddleware(authHandler)

	fmt.Println("API Gateway running on port 8080...")
	if err := http.ListenAndServe(":8080", finalHandler); err != nil {
		log.Fatal(err)
	}
}
