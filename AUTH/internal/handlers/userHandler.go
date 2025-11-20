package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	proto "PROJEKAT/COMMON/stakeholders/proto"
	"auth/internal/models"
	"auth/internal/services"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	h.SetProfile(*user) // kreiraj profil u stakeholders servisu

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// mapUserRole konvertuje string iz baze u protobuf enum
func mapUserRole(r models.Role) proto.Role {
	switch r {
	case models.RoleGuide:
		return proto.Role_ROLE_GUIDE
	case models.RoleTourist:
		return proto.Role_ROLE_TOURIST
	case models.RoleAdmin:
		return proto.Role_ROLE_ADMIN
	default:
		return proto.Role_ROLE_UNKNOWN
	}
}

func (h *UserHandler) SetProfile(user models.User) {
	conn, err := grpc.Dial("stakeholders:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewStakeholdersClient(conn)

	req := &proto.CreateProfileRequest{
		UserId: user.Id.String(), // generisano iz protobuf -> UserId
		Role:   mapUserRole(user.Role),
	}

	resp, err := client.CreateProfile(context.Background(), req)
	if err != nil {
		fmt.Println("Error creating profile:", err)
		return
	}

	fmt.Printf("Profile created: id=%s, already_existed=%v\n", resp.Id, resp.AlreadyExisted)
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
