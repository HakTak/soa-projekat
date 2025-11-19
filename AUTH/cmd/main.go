package main

import (
	"fmt"
	"log"
	"net/http"

	"auth/internal/database"
	"auth/internal/handlers"
	"auth/internal/middleware"
	"auth/internal/repositories"
	"auth/internal/services"
)

func main() {
	// 1. konekcija na bazu
	db := database.Connect()

	// 2. repo + service + handler
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	jwtService := services.NewJWTService(
		"my_super_duper_secret_mega_gg_key_123", //secretKey
		"AUTH_SERVICE",                          //issuer
	)
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	userHandler := handlers.NewUserHandler(userService, jwtService)

	// 3. rute
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{\"status\":\"AUTH is ok\"}")
	})

	http.HandleFunc("/register", userHandler.Register)
	http.HandleFunc("/login", userHandler.Login)

	http.HandleFunc("/users", userHandler.GetAll)
	http.HandleFunc("/admin/users", authMiddleware.AdminOnly(userHandler.GetUsersForAdmin))

	fmt.Println("Auth running on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
