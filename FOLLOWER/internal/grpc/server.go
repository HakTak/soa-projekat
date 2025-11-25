package grpc

import (
	"FOLLOWER/internal/service"
	"context"

	pb "FOLLOWER/common/genproto"
)

type FollowerGrpcServer struct {
	pb.UnimplementedFollowerServiceServer
	svc service.FollowService
}

func NewFollowerGrpcServer(svc service.FollowService) *FollowerGrpcServer {
	return &FollowerGrpcServer{svc: svc}
}

func (s *FollowerGrpcServer) IsFollowing(ctx context.Context, req *pb.IsFollowingRequest) (*pb.IsFollowingResponse, error) {
	isFollowing, err := s.svc.IsFollowing(ctx, req.FollowerId, req.FolloweeId)
	if err != nil {
		return nil, err
	}

	return &pb.IsFollowingResponse{
		IsFollowing: isFollowing,
	}, nil
}
