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

	if err := h.service.CreateTour(&t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(t)
}

func (h *TourHandler) GetTour(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	tour, err := h.service.GetTour(uint(id))
	if err != nil {
		http.Error(w, "tour not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(tour)
}

func (h *TourHandler) GetAllTours(w http.ResponseWriter, r *http.Request) {
	tours, err := h.service.GetAllTours()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tours)
}

func (h *TourHandler) DeleteTour(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	if err := h.service.DeleteTour(uint(id)); err != nil {
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
