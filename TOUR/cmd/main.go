package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"tour/internal/database"
	handlers "tour/internal/handler"
	"tour/internal/repository"
	"tour/internal/service"
)

func main() {
	// 1️⃣ Konekcija na bazu
	db := database.Connect() // tu u Connect() treba da bude gorm.Open sa DSN-om

	// 2️⃣ Repo + Service + Handler
	tourRepo := repository.NewGormTourRepo(db)
	tourService := service.NewTourService(tourRepo)
	tourHandler := handlers.NewTourHandler(tourService)

	// 3️⃣ Health check ruta
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{\"status\":\"TOUR service is ok\"}")
	})

	// 4️⃣ Tour rute
	//http.HandleFunc("/tours", tourHandler.GetAll)         // GET sve ture
	http.HandleFunc("/tours/create", tourHandler.Create)                  // POST nova tura
	http.HandleFunc("/tours/getByAuthorID", tourHandler.GetAllByAuthorID) // GET jedna tura po ID
	//http.HandleFunc("/tours/update/", tourHandler.Update) // PUT update ture
	//http.HandleFunc("/tours/delete/", tourHandler.Delete) // DELETE tura

	// 5️⃣ Pokretanje servera
	srv := &http.Server{
		Addr:         ":8081",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("Tour service running on port 8081...")
	log.Fatal(srv.ListenAndServe())
}
