package main

import (
	"log"
	"net"

	"SHOPPING-CART/internal/database"
	"SHOPPING-CART/internal/handlers"
	"SHOPPING-CART/internal/repository"
	"SHOPPING-CART/internal/service"

	// ✅ IMPORT 1: Shopping Cart Proto (Server)
	pbCart "PROJEKAT/COMMON/shopping-cart/proto"

	// ✅ IMPORT 2: Tour Proto (Client) - Using the COMMON module
	pbTour "PROJEKAT/COMMON/tour/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 1. Database Connection
	db := database.Connect()

	// 2. Repository Initialization
	repo := repository.NewCartRepository(db)

	// ========================================================================
	// 3. Tour Service Client Setup
	// ========================================================================

	// Connect to the Tour container
	tourConn, err := grpc.Dial("tour:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Tour Service: %v", err)
	}
	defer tourConn.Close()

	// Create the Client instance using the COMMON proto definition
	tourClient := pbTour.NewTourServiceClient(tourConn)

	// ========================================================================
	// 4. Service & Handler Initialization
	// ========================================================================

	svc := service.NewShoppingCartService(repo, tourClient)
	cartHandler := handlers.NewShoppingCartHandler(svc)

	// ========================================================================
	// 5. Start gRPC Server
	// ========================================================================

	// Ensure this port matches what you registered in API Gateway (9092)
	lis, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pbCart.RegisterShoppingCartServiceServer(grpcServer, cartHandler)
	reflection.Register(grpcServer)

	log.Println("Shopping Cart Microservice (gRPC) started on port 9092")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
