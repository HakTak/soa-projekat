package main

import (
	"log"
	"net"

	// Update "SHOPPING-CART" to match your actual go.mod module name
	"SHOPPING-CART/internal/database"
	"SHOPPING-CART/internal/handlers"
	"SHOPPING-CART/internal/repository"
	"SHOPPING-CART/internal/service"

	// Import the generated Protocol Buffers
	// 1. The Cart Proto (We are the Server for this)
	pbCart "PROJEKAT/COMMON/shopping-cart/proto"

	// 2. The Tour Proto (We are the Client for this)
	// Make sure this path points to where your Tour proto is generated
	pbTour "SHOPPING-CART/common/genproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 1. Database Connection
	// Using the Connect function we created in internal/database/database.go
	db := database.Connect()

	// 2. Repository Initialization
	repo := repository.NewCartRepository(db)

	// ========================================================================
	// 3. Tour Service Client Setup (Crucial Step)
	// The Shopping Cart needs to talk to the Tour Service to validate tours.
	// ========================================================================

	// "tour:50052" matches the service name and gRPC port in docker-compose
	tourConn, err := grpc.Dial("tour:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Tour Service: %v", err)
	}
	defer tourConn.Close()

	// Create the Client instance
	tourClient := pbTour.NewTourServiceClient(tourConn)

	// ========================================================================
	// 4. Service & Handler Initialization
	// ========================================================================

	// Inject both the Repository (DB) and the Tour Client (Microservice)
	svc := service.NewShoppingCartService(repo, tourClient)

	// Create the gRPC Handler
	cartHandler := handlers.NewShoppingCartHandler(svc)

	// ========================================================================
	// 5. Start gRPC Server
	// ========================================================================

	// Listen on port 50053 (Must match docker-compose for shopping-cart)
	lis, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register the Cart Handler
	pbCart.RegisterShoppingCartServiceServer(grpcServer, cartHandler)

	// Enable Reflection (Optional, but good for testing with Postman/gRPCurl)
	reflection.Register(grpcServer)

	log.Println("Shopping Cart Microservice (gRPC) started on port 50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
