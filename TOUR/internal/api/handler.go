package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tour-service/internal/model"
	"tour-service/internal/service"

	"github.com/go-chi/chi/v5"
)

type TourHandler struct {
	service *service.TourService
}

func NewTourHandler(s *service.TourService) *TourHandler {
	return &TourHandler{service: s}
}

func (h *TourHandler) CreateTour(w http.ResponseWriter, r *http.Request) {
	var t model.Tour
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t.Status = "DRAFT"

	if err := h.service.CreateTour(&t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(t)
}

func (h *TourHandler) GetTour(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	tour, err := h.service.GetTour(id)
	if err != nil {
		http.Error(w, "tour not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(tour)
}

func (h *TourHandler) GetAllTours(w http.ResponseWriter, r *http.Request) {
	// Dozvoli pristup sa Angulara (localhost:4200) ili stavi "*" za sve
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Ako browser salje "preflight" OPTIONS zahtev, odmah vrati OK
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	tours, err := h.service.GetAllTours()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tours)
}

func (h *TourHandler) DeleteTour(w http.ResponseWriter, r *http.Request) {
	/*
		idStr := chi.URLParam(r, "id")
		id, _ := strconv.ParseUint(idStr, 10, 64)*/

	id := chi.URLParam(r, "id")

	if err := h.service.DeleteTour(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TourHandler) UpdateTour(w http.ResponseWriter, r *http.Request) {
	var t model.Tour
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.UpdateTour(&t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(t)
}

func (h *TourHandler) GetToursByUser(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "userId")
	userId, err := strconv.ParseUint(userIdStr, 10, 64) // ne znam sto se parisra ali ajde

	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}
	tours, err := h.service.GetToursByUser(uint(userId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tours)
}

//****** Review handlers ******//

func (h *TourHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var rev model.Review
	if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.CreateReview(&rev); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rev)
}

func (h *TourHandler) GetReviewsByTour(w http.ResponseWriter, r *http.Request) {

	// Dozvoli pristup sa Angulara (localhost:4200) ili stavi "*" za sve
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Ako browser salje "preflight" OPTIONS zahtev, odmah vrati OK
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	tourID := chi.URLParam(r, "tourId")

	reviews, err := h.service.GetReviewsByTour(tourID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(reviews)
}

func (h *TourHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteReview(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Review deleted successfully"})

}

func (h *TourHandler) GetAllReviews(w http.ResponseWriter, r *http.Request) {
	var err error
	var reviews []model.Review
	reviews, err = h.service.GetAllReviews()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(reviews)
}
