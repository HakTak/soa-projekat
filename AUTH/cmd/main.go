package main

import (
	"fmt"
	"log"
	"net"

	"auth/internal/database"
	"auth/internal/handlers"
	"auth/internal/repositories"
	"auth/internal/services"

	pbAuth "PROJEKAT/COMMON/auth/proto"
	commonMiddleware "PROJEKAT/COMMON/middleware"
	pbStakeholders "PROJEKAT/COMMON/stakeholders/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	db := database.Connect()

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	jwtService := services.NewJWTService("my_super_duper_secret_mega_gg_key_123", "AUTH_SERVICE")

	// --- NOVO: Konekcija ka Stakeholders servisu ---
	// Ovo radimo samo jednom pri startovanju aplikacije
	stakeholdersConn, err := grpc.NewClient("stakeholders:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to stakeholders: %v", err)
	}
	defer stakeholdersConn.Close()

	// Pravimo klijenta
	stakeholdersClient := pbStakeholders.NewStakeholdersServiceClient(stakeholdersConn)
	// -----------------------------------------------

	// Prosledjujemo klijenta u handler
	userHandler := handlers.NewUserHandler(userService, jwtService, stakeholdersClient)

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(commonMiddleware.MetadataExtractorInterceptor),
	)

	pbAuth.RegisterAuthServiceServer(grpcServer, userHandler)

	fmt.Println("Auth gRPC Service running on port 8082...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
