package handlers

import (
	"auth/internal/models"
	"auth/internal/services"
	"encoding/json"
	"net/http"
	"strings"
)

type UserHandler struct {
	userService *services.UserService
	jwtService  *services.JWTService
}

func NewUserHandler(s *services.UserService, j *services.JWTService) *UserHandler {
	return &UserHandler{s, j}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

	// 1. PROVERA — DA LI JE KORISNIK VEĆ ULOGOVAN
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {

		// izvuci token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := h.jwtService.ValidateToken(tokenStr)

		if err == nil && token.Valid {
			http.Error(w, "You are already logged in — cannot register a new account", http.StatusForbidden)
			return
		}
	}

	// 2. PARSE BODY
	var userRegister models.UserRegisterDTO
	err := json.NewDecoder(r.Body).Decode(&userRegister)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3. VALIDACIJA ROLE — DOZVOLJENO SAMO GUIDE ILI TOURIST
	if userRegister.Role == models.RoleAdmin {
		http.Error(w, "You cannot register as ADMIN", http.StatusForbidden)
		return
	}

	if userRegister.Role != models.RoleGuide && userRegister.Role != models.RoleTourist {
		http.Error(w, "Invalid role — allowed: GUIDE or TOURIST", http.StatusBadRequest)
		return
	}

	// 4. POZIV SERVISA
	user, err := h.userService.Register(userRegister)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Password = "*********" // ne šaljemo password

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq models.LoginRequestDTO
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(loginReq.Username, loginReq.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := h.jwtService.GenerateToken(user.Id.String(), string(user.Role))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := models.LoginResponseDTO{
		Token: token,
		User: models.UserNoPassDTO{
			Id:       user.Id.String(),
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Blocked:  user.Blocked,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// // ukloni password iz liste
	// for i := range users {
	// 	users[i].Password = ""
	// }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUsersForAdmin(w http.ResponseWriter, r *http.Request) {
	nonAdminUsers, err := h.userService.GetUsersForAdmin()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nonAdminUsers)
}
