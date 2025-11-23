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
	db.AutoMigrate(&model.Tour{}, &model.Keypoint{}, &model.Review{})

	// Repository -> Service -> Handler
	repo := repository.NewTourRepository(db)
	reviewRepo := repository.CreateReviewRepository(db)
	svc := service.NewTourService(repo, reviewRepo)
	handler := api.NewTourHandler(svc)

	r := chi.NewRouter()

	// ---------------------------------------------------------
	// CORS MIDDLEWARE - POCETAK
	// Ovo omogucava Angular aplikaciji (localhost:4200) da komunicira sa ovim servisom
	// ---------------------------------------------------------
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Dozvoli zahteve sa Angular porta
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")

			// Dozvoli metode koje koristis u rutama (ukljucujuci PATCH i DELETE)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")

			// Dozvoli standardne hedere
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")

			// Dozvoli kredencijale ako budu potrebni
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Ako browser salje "OPTIONS" zahtev (preflight check), odmah vrati OK
			// i ne salji zahtev dalje ka handlerima
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})
	// ---------------------------------------------------------
	// CORS MIDDLEWARE - KRAJ
	// ---------------------------------------------------------

	r.Post("/tour", handler.CreateTour)
	r.Get("/tour/{id}", handler.GetTour)
	r.Get("/tours", handler.GetAllTours)
	r.Get("/tours/myTours/{userId}", handler.GetToursByUser)
	r.Delete("/tour/{id}", handler.DeleteTour)
	r.Patch("/tour/update", handler.UpdateTour)

	r.Post("/review/create", handler.CreateReview)
	r.Get("/review/tour/{tourId}", handler.GetReviewsByTour)
	r.Get("/review/getAll", handler.GetAllReviews)
	r.Delete("/review/delete/{id}", handler.DeleteReview)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Println("Tour service running on port", port)
	http.ListenAndServe(":"+port, r)
}
