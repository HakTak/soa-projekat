package grpc

import (
	"context"
	"errors"
	"fmt"
	"stakeholders/internal/model"
	"stakeholders/internal/service"

	pb "stakeholders/common/genproto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedStakeholdersServer
	svc service.ProfileService
}

func NewServer(svc service.ProfileService) *Server {
	return &Server{svc: svc}
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
		return "", fmt.Errorf("unkown role: %v", r)
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

func (s *Server) CreateProfile(ctx context.Context, req *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id required")
	}
	if req.Role == pb.Role_ROLE_UNKNOWN {
		return nil, status.Error(codes.InvalidArgument, "valid role required")
	}

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

	created, err := s.svc.CreateProfile(ctx, p)
	if err != nil {
		if errors.Is(err, service.ErrProfileExists) {
			existing, gerr := s.svc.GetProfile(ctx, req.UserId)
			if gerr != nil {
				return nil, status.Error(codes.Internal, "profile exists but failed to fetch")
			}
			return &pb.CreateProfileResponse{Id: existing.ID, AlreadyExisted: true}, nil
		}
		return nil, status.Errorf(codes.Internal, "create profile failed do to %v", err)
	}
	return &pb.CreateProfileResponse{Id: created.ID, AlreadyExisted: false}, nil
}

func (s *Server) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id required")
	}
	p, err := s.svc.GetProfile(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
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
	}, nil
}
