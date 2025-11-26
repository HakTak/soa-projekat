package handlers

import (
	"context"
	"fmt"

	// Importujes zajednicke pakete
	pb "PROJEKAT/COMMON/stakeholders/proto"
	"PROJEKAT/COMMON/utils"

	"stakeholders/internal/model"
	"stakeholders/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProfileHandler struct {
	pb.UnimplementedStakeholdersServiceServer
	svc service.ProfileService
}

func NewProfileHandler(svc service.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

// ==========================================
// 1. INTERNE METODE (Poziva Auth Servis)
// ==========================================

func (h *ProfileHandler) CreateProfile(ctx context.Context, req *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id required")
	}

	// Koristimo helper za konverziju role
	role, err := protoRoleToModel(req.Role)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid role: %v", err)
	}

	p := &model.Profile{
		UserID:         req.UserId,
		FirstName:      "",
		LastName:       "",
		ProfilePicture: "",
		Biography:      "",
		Motto:          "",
		Role:           role,
		IsBlocked:      false,
	}

	created, err := h.svc.CreateProfile(ctx, p)
	if err != nil {
		// Idempotentnost: Ako postoji, vratimo uspeh i ID postojeceg
		if err == service.ErrProfileExists {
			existing, _ := h.svc.GetProfile(ctx, req.UserId)
			return &pb.CreateProfileResponse{Id: existing.ID, AlreadyExisted: true}, nil
		}
		return nil, status.Errorf(codes.Internal, "create profile failed: %v", err)
	}
	return &pb.CreateProfileResponse{Id: created.ID, AlreadyExisted: false}, nil
}

// ==========================================
// 2. JAVNE METODE (Poziva API Gateway)
// ==========================================

func (h *ProfileHandler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	// Logika autorizacije:
	// Korisnik moze da vidi profil ako je TO NJEGOV profil ILI ako je ADMIN.

	claims := utils.ClaimsFromContext(ctx)

	// Ako nema claims (npr. nije doslo preko Gateway-a ili nema tokena),
	// mozda zelimo da dozvolimo javni pristup (ako je profil javan)?
	// Ali po tvojoj logici, mora biti vlasnik ili admin.

	if claims == nil {
		return nil, status.Error(codes.Unauthenticated, "Authentication required")
	}

	requesterID := claims["id"].(string)
	requesterRole := claims["role"].(string)

	// Ako trazilac nije vlasnik I nije admin -> Forbidden
	if requesterID != req.UserId && requesterRole != "ADMIN" {
		return nil, status.Error(codes.PermissionDenied, "Access denied")
	}

	p, err := h.svc.GetProfile(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Profile not found")
	}
	return mapProfileToProto(p), nil
}

func (h *ProfileHandler) GetMyProfile(ctx context.Context, _ *emptypb.Empty) (*pb.GetProfileResponse, error) {
	claims, err := h.getClaimsSafe(ctx)
	if err != nil {
		return nil, err
	}

	myID := claims["id"].(string)

	// 2. Pozovi servis
	p, err := h.svc.GetProfile(ctx, myID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Profile not found")
	}

	// 3. Proveri da li je nil (za svaki slucaj)
	if p == nil {
		return nil, status.Error(codes.NotFound, "Profile data is empty")
	}

	return mapProfileToProto(p), nil
}

func (h *ProfileHandler) UpdateMyProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.GetProfileResponse, error) {
	claims, err := h.getClaimsSafe(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Safe cast (we know "id" exists because getClaimsSafe checked it)
	myID := claims["id"].(string)

	existing, err := h.svc.GetProfile(ctx, myID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Profile not found")
	}

	if existing == nil {
		return nil, status.Error(codes.NotFound, "Profile data is empty")
	}

	// Mapiranje polja (Update logic)
	if req.FirstName != "" {
		existing.FirstName = req.FirstName
	}
	if req.LastName != "" {
		existing.LastName = req.LastName
	}
	if req.ProfilePicture != "" {
		existing.ProfilePicture = req.ProfilePicture
	}
	if req.Biography != "" {
		existing.Biography = req.Biography
	}
	if req.Motto != "" {
		existing.Motto = req.Motto
	}

	updated, err := h.svc.UpdateProfile(ctx, existing)
	if err != nil {
		return nil, status.Error(codes.Internal, "Update failed")
	}
	return mapProfileToProto(updated), nil
}

func (h *ProfileHandler) BlockUser(ctx context.Context, req *pb.BlockUserRequest) (*emptypb.Empty, error) {
	// Samo ADMIN moze da blokira
	if err := utils.Authorize(ctx, "ADMIN"); err != nil {
		return nil, err
	}

	claims := utils.ClaimsFromContext(ctx)
	adminID := claims["id"].(string)

	err := h.svc.BlockUser(ctx, req.UserId, adminID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Block failed: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// ==========================================
// 3. POMOCNE FUNKCIJE (Maperi)
// ==========================================

func mapProfileToProto(p *model.Profile) *pb.GetProfileResponse {
	if p == nil {
		return &pb.GetProfileResponse{} // Return empty object or handle differently
	}
	return &pb.GetProfileResponse{
		Id:             p.ID,
		UserId:         p.UserID,
		FirstName:      p.FirstName,
		LastName:       p.LastName,
		Role:           modelRoleToProto(p.Role),
		ProfilePicture: p.ProfilePicture,
		Biography:      p.Biography,
		Motto:          p.Motto,
		IsBlocked:      p.IsBlocked,
	}
}

func modelRoleToProto(r model.Role) pb.Role {
	switch r {
	case model.RoleAdmin:
		return pb.Role_ROLE_ADMIN
	case model.RoleGuide:
		return pb.Role_ROLE_GUIDE
	case model.RoleTourist:
		return pb.Role_ROLE_TOURIST
	default:
		return pb.Role_ROLE_UNKNOWN
	}
}

func protoRoleToModel(r pb.Role) (model.Role, error) {
	switch r {
	case pb.Role_ROLE_GUIDE:
		return model.RoleGuide, nil
	case pb.Role_ROLE_ADMIN:
		return model.RoleAdmin, nil
	case pb.Role_ROLE_TOURIST:
		return model.RoleTourist, nil
	default:
		return "", fmt.Errorf("unknown role: %v", r)
	}
}

func (h *ProfileHandler) getClaimsSafe(ctx context.Context) (map[string]interface{}, error) {
	// 1. Try getting claims from the shared utility
	claims := utils.ClaimsFromContext(ctx)

	// If utils returned a map, ensure it has the ID
	if claims != nil {
		if _, ok := claims["id"].(string); ok {
			return claims, nil
		}
	}

	// 2. FALLBACK: Read directly from gRPC Metadata
	// (This handles cases where the Interceptor might have failed or key names mismatch)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "No metadata provided")
	}

	// Create a temporary claims map
	newClaims := make(map[string]interface{})

	// Gateway sends "user-id" or "x-user-id"
	if ids := md.Get("user-id"); len(ids) > 0 {
		newClaims["id"] = ids[0]
	} else if ids := md.Get("x-user-id"); len(ids) > 0 {
		newClaims["id"] = ids[0]
	}

	// Gateway sends "user-role" or "x-user-role"
	if roles := md.Get("user-role"); len(roles) > 0 {
		newClaims["role"] = roles[0]
	} else if roles := md.Get("x-user-role"); len(roles) > 0 {
		newClaims["role"] = roles[0]
	}

	// Check if we found the ID
	if _, ok := newClaims["id"]; !ok {
		return nil, status.Error(codes.Unauthenticated, "User ID not found in request")
	}

	return newClaims, nil
}
