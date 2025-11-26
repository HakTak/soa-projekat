package main

import (
	"context"
	"log"
	"net"

	"FOLLOWER/internal/config"
	"FOLLOWER/internal/database"
	"FOLLOWER/internal/handlers"
	"FOLLOWER/internal/repository"
	"FOLLOWER/internal/service"

	pb "PROJEKAT/COMMON/follower/proto"
	commonMiddleware "PROJEKAT/COMMON/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.LoadConfig()

	driver, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Cannot connect to Neo4j: %v", err)
	}
	defer driver.Close(context.Background())

	repo := repository.NewfollowRepository(driver)

	log.Println("Applying Neo4j constraints...")
	if err := repo.EnsureConstraints(context.Background()); err != nil {
		log.Printf("Warning: Could not ensure constraints: %v", err)
	} else {
		log.Println("Constraints applied.")
	}

	svc := service.NewFollowService(repo)

	followerHandler := handlers.NewFollowerHandler(svc)

	//start grpc server port 9090
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(commonMiddleware.MetadataExtractorInterceptor),
	)
	pb.RegisterFollowerServiceServer(grpcServer, followerHandler)

	reflection.Register(grpcServer)

	log.Println("Follower Microservice (gRPC) started on port 9090")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
