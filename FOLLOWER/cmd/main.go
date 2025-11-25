package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	pb "FOLLOWER/common/genproto"
	"FOLLOWER/internal/config"
	mygrpc "FOLLOWER/internal/grpc"
	myhttp "FOLLOWER/internal/http" // Alias to avoid collision with net/http
	"FOLLOWER/internal/repository"
	"FOLLOWER/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.LoadConfig() // (Reuse the config code provided in previous steps)

	// 2. Database
	driver, err := connectToNeo4j(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Close(context.Background())

	// 3. Dependency Injection
	repo := repository.NewfollowRepository(driver)

	log.Println("Applying Neo4j constraints...")
	if err := repo.EnsureConstraints(context.Background()); err != nil {
		log.Fatalf("❌ Failed to ensure constraints: %v", err)
	}
	log.Println("✅ Constraints applied.")

	svc := service.NewFollowService(repo)
	handler := myhttp.NewHandler(svc)

	// gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":9090") // gRPC Port
		if err != nil {
			log.Fatalf("failed to listen for gRPC: %v", err)
		}

		grpcServer := grpc.NewServer()

		// Register the service
		followerServer := mygrpc.NewFollowerGrpcServer(svc)
		pb.RegisterFollowerServiceServer(grpcServer, followerServer)

		// Enable reflection (optional, helps with testing using Postman/Evans)
		reflection.Register(grpcServer)

		log.Println("✅ gRPC Server started on port 9090")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// 4. Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(myhttp.AuthMiddleware)

	handler.RegisterRoutes(r)

	// 5. Run
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	log.Printf("Starting Social Service on port %s", port)
	http.ListenAndServe(":"+port, r)
}

// connectToNeo4j logic (retry loop) goes here...
func connectToNeo4j(cfg *config.Config) (neo4j.DriverWithContext, error) {
	// ... (Same retry logic as provided in previous answer)
	// For brevity, assuming simple connection:
	return neo4j.NewDriverWithContext(cfg.Neo4jURI, neo4j.BasicAuth(cfg.Neo4jUser, cfg.Neo4jPassword, ""))
}
