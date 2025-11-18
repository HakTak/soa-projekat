package main

import (
    "fmt"
    "log"
    "net/http"

    "auth/internal/database"
    "auth/internal/repositories"
    "auth/internal/services"
    "auth/internal/handlers"
)

func main() {
    // 1. konekcija na bazu
    db := database.Connect()

    // 2. repo + service + handler
    userRepo := repositories.NewUserRepository(db)
    userService := services.NewUserService(userRepo)
    userHandler := handlers.NewUserHandler(userService)

    // 3. rute
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "{\"status\":\"AUTH is ok\"}")
    })

    http.HandleFunc("/register", userHandler.Register)
    http.HandleFunc("/users", userHandler.GetAll)

    fmt.Println("Auth running on port 8082...")
    log.Fatal(http.ListenAndServe(":8082", nil))
}
