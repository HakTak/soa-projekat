package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"tour/internal/model"
	"tour/internal/service"
)

type TourHandler struct {
	service *service.TourService
}

func NewTourHandler(s *service.TourService) *TourHandler {
	return &TourHandler{s}
}

// Request payload za kreiranje ture
type CreateTourRequest struct {
	AuthorID    string               `json:"author_id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Dificulty   model.DificultyLevel `json:"dificulty_level"`
	Status      model.Status         `json:"status"`
	Tags        []model.Tag          `json:"tags"`
}

// Create nova tura
func (h *TourHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTourRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	println("AuthorID: ", req.AuthorID)
	tour, err := h.service.Create(r.Context(), req.AuthorID, req.Name, req.Description, req.Dificulty, req.Status, req.Tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tour)
}

// GetAllByAuthorID vraća sve ture jednog autora
func (h *TourHandler) GetAllByAuthorID(w http.ResponseWriter, r *http.Request) {

	authorID := strings.TrimSpace(r.URL.Query().Get("author_id"))

	if authorID == "" {
		http.Error(w, "author_id query param required", http.StatusBadRequest)
		return
	}

	tours, err := h.service.GetAllByAuthorID(r.Context(), authorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tours)
}
