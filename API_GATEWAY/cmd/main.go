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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 1. GRPC GATEWAY MUX
	gwmux := runtime.NewServeMux(
	// runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
	// 	switch key {
	// 	case "Grpc-Metadata-User-Id":
	// 		return "x-user-id", true
	// 	case "Grpc-Metadata-User-Role":
	// 		return "x-user-role", true
	// 	}
	// 	return runtime.DefaultHeaderMatcher(key)
	// }),
	)
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
