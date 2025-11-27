package main

import (
	"log"
	"net"
	"os"
	"tour-service/internal/api"
	"tour-service/internal/model"
	"tour-service/internal/repository"
	"tour-service/internal/service"

	commonMiddleware "PROJEKAT/COMMON/middleware"
	pb "PROJEKAT/COMMON/tour/proto"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
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

	// Auto migrate tables
	db.AutoMigrate(&model.Tour{}, &model.Keypoint{}, &model.Review{}, &model.RouteOption{})

	// Repository -> Service -> Handler
	repo := repository.NewTourRepository(db)
	reviewRepo := repository.CreateReviewRepository(db)
	svc := service.NewTourService(repo, reviewRepo)
	handler := api.NewTourGRPCServer(svc)
	// --- gRPC server ---
	lis, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Ubacujemo zajednicki interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(commonMiddleware.MetadataExtractorInterceptor),
	)

	pb.RegisterTourServiceServer(grpcServer, handler)
	// <- handler must implement pb.TourServiceServer

	log.Println("Tour gRPC service running on port 8083")
	grpcServer.Serve(lis)
}
