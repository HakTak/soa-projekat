package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"tour/internal/database"
	"tour/internal/handler"
	"tour/internal/repository"
	"tour/internal/service"
)

func main() {
	// 1️⃣ Konekcija na bazu
	db := database.Connect()

	// 2️⃣ Repo + Service + Handler za Tour
	tourRepo := repository.NewGormTourRepo(db)
	tourService := service.NewTourService(tourRepo)
	tourHandler := handler.NewTourHandler(tourService)

	// 2️⃣ Repo + Service + Handler za Comment
	commentRepo := repository.NewCommentRepository(db)
	commentService := service.NewCommentService(commentRepo, tourRepo)
	commentHandler := handler.NewCommentHandler(commentService)

	// 3️⃣ Health check ruta
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{\"status\":\"TOUR service is ok\"}")
	})

	// 4️⃣ Tour rute
	http.HandleFunc("/tours/create", tourHandler.Create)
	http.HandleFunc("/tours/getByAuthorID", tourHandler.GetAllByAuthorID)

	// 5️⃣ Comment rute
	http.HandleFunc("/comments/create", commentHandler.Create)           // POST nova recenzija
	http.HandleFunc("/comments/getByTourID", commentHandler.GetComments) // GET recenzije za jednu turu

	// 6️⃣ Pokretanje servera
	srv := &http.Server{
		Addr:         ":8094",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("Tour service running on port 8094...")
	log.Fatal(srv.ListenAndServe())
}
