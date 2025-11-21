package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"tour/internal/service"
)

type CommentHandler struct {
	service *service.CommentService
}

func NewCommentHandler(s *service.CommentService) *CommentHandler {
	return &CommentHandler{s}
}

type CreateCommentRequest struct {
	TourID   string    `json:"tour_id"`
	UserID   string    `json:"user_id"`
	UserName string    `json:"user_name"`
	Content  string    `json:"content"`
	ImageURL string    `json:"image_url"`
	Rating   int       `json:"rating"`
	TourDate time.Time `json:"tour_date"`
}

// Create novi komentar
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	comment, err := h.service.Create(r.Context(), req.TourID, req.UserID, req.UserName, req.Content, req.ImageURL, req.Rating, req.TourDate)
	if err != nil {

		if err.Error() == "tour does not exist" {
			http.Error(w, `{"error":"tour does not exist"}`, http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comment)
}

// GetComments vraća sve komentare za jednu turu
func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	tourID := r.URL.Query().Get("tour_id")
	if tourID == "" {
		http.Error(w, "tour_id query param required", http.StatusBadRequest)
		return
	}

	comments, err := h.service.GetByTourID(r.Context(), tourID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
}
