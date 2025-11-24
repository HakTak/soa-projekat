package handlers

import (
	"context"
	"fmt"

	pbAuth "PROJEKAT/COMMON/auth/proto"
	pbStakeholders "PROJEKAT/COMMON/stakeholders/proto"

	// IMPORTUJEMO UTILS
	"PROJEKAT/COMMON/utils"

	"auth/internal/models"
	"auth/internal/services"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	pbAuth.UnimplementedAuthServiceServer
	userService        *services.UserService
	jwtService         *services.JWTService
	stakeholdersClient pbStakeholders.StakeholdersServiceClient // <-- NOVO: Klijent je ovde
}

// Primamo klijenta u konstruktoru
func NewUserHandler(s *services.UserService, j *services.JWTService, sc pbStakeholders.StakeholdersServiceClient) *UserHandler {
	return &UserHandler{
		userService:        s,
		jwtService:         j,
		stakeholdersClient: sc,
	}
}

// REGISTER
func (h *UserHandler) Register(ctx context.Context, req *pbAuth.RegisterRequest) (*pbAuth.RegisterResponse, error) {
	reqRole := models.Role(req.Role)

	if reqRole == models.RoleAdmin {
		return nil, status.Error(codes.PermissionDenied, "Cannot register as ADMIN")
	}
	if reqRole != models.RoleGuide && reqRole != models.RoleTourist {
		return nil, status.Error(codes.InvalidArgument, "Role must be GUIDE or TOURIST")
	}

	dto := models.UserRegisterDTO{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     reqRole,
	}

	user, err := h.userService.Register(dto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Registration failed: %v", err)
	}

	// Sada samo pozivamo metodu, nema konektovanja
	h.createStakeholderProfile(ctx, *user)

	return &pbAuth.RegisterResponse{
		Id:       user.Id.String(),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// LOGIN
func (h *UserHandler) Login(ctx context.Context, req *pbAuth.LoginRequest) (*pbAuth.LoginResponse, error) {
	user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}

	token, err := h.jwtService.GenerateToken(
		user.Id.String(),
		string(user.Role),
		user.Username,
		user.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "Token generation failed")
	}

	return &pbAuth.LoginResponse{Token: token}, nil
}

// ADMIN ONLY
func (h *UserHandler) GetUsersForAdmin(ctx context.Context, _ *emptypb.Empty) (*pbAuth.GetAllUsersResponse, error) {

	// 1. KORISTIMO NOVU UTILS FUNKCIJU
	// Ovo proverava da li korisnik u kontekstu ima ulogu "ADMIN"
	if err := utils.Authorize(ctx, "ADMIN"); err != nil {
		return nil, err
	}

	// 2. Ako prodje, nastavljamo logiku
	users, err := h.userService.GetUsersForAdmin()
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to fetch users")
	}
	return h.mapUsersToProto(users), nil
}

// --- POMOCNE FUNKCIJE ---

func (h *UserHandler) mapUsersToProto(users []models.UserNoPassDTO) *pbAuth.GetAllUsersResponse {
	// Pazi: ovde sam vratio []models.User jer to vraca tvoj servis,
	// ako si menjao DTO, prilagodi ovo.
	var protoUsers []*pbAuth.UserResponse
	for _, u := range users {
		protoUsers = append(protoUsers, &pbAuth.UserResponse{
			Id:       u.Id,
			Username: u.Username,
			Email:    u.Email,
			Role:     string(u.Role),
			Blocked:  u.Blocked,
		})
	}
	return &pbAuth.GetAllUsersResponse{Users: protoUsers}
}

// REFAKTORISANO: Koristi postojeceg klijenta
func (h *UserHandler) createStakeholderProfile(ctx context.Context, user models.User) {
	var protoRole pbStakeholders.Role
	switch user.Role {
	case models.RoleGuide:
		protoRole = pbStakeholders.Role_ROLE_GUIDE
	case models.RoleTourist:
		protoRole = pbStakeholders.Role_ROLE_TOURIST
	default:
		protoRole = pbStakeholders.Role_ROLE_UNKNOWN
	}

	// Koristimo h.stakeholdersClient koji je vec povezan
	_, err := h.stakeholdersClient.CreateProfile(ctx, &pbStakeholders.CreateProfileRequest{
		UserId: user.Id.String(),
		Role:   protoRole,
	})

	if err != nil {
		fmt.Printf("Warning: Profile creation failed: %v\n", err)
	} else {
		fmt.Println("Profile created in Stakeholders service")
	}
}
