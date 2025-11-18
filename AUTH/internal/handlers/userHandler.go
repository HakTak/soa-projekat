package handlers

import (
	"auth/internal/models"
	"auth/internal/services"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{s}
}

type RegisterRequest struct {
	Username string       `json:"username"`
	Email    string       `json:"email"`
	Password string       `json:"password"`
	Role     models.Role  `json:"role"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	json.NewDecoder(r.Body).Decode(&req)

	user, err := h.service.Register(req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// user.Password = "" // ne smeš slati lozinku
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// // ukloni password iz liste
	// for i := range users {
	// 	users[i].Password = ""
	// }

	json.NewEncoder(w).Encode(users)
}
