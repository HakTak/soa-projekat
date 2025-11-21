package main

import (
	"log"
	"net/http"
	"os"
	"tour-service/internal/api"
	"tour-service/internal/model"
	"tour-service/internal/repository"
	"tour-service/internal/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/go-chi/chi/v5"
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
	db.AutoMigrate(&model.Tour{}, &model.Keypoint{})

	// Repository -> Service -> Handler
	repo := repository.NewTourRepository(db)
	svc := service.NewTourService(repo)
	handler := api.NewTourHandler(svc)

	r := chi.NewRouter()

	r.Post("/tour", handler.CreateTour)
	r.Get("/tour/{id}", handler.GetTour)
	r.Get("/tours", handler.GetAllTours)
	r.Delete("/tour/{id}", handler.DeleteTour)
	r.Patch("/tour/update", handler.UpdateTour)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Println("Tour service running on port", port)
	http.ListenAndServe(":"+port, r)
}
