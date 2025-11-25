package handlers

import (
	"context"

	// ✅ Import the generated code from COMMON
	"FOLLOWER/internal/service"
	pb "PROJEKAT/COMMON/follower/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type FollowerHandler struct {
	pb.UnimplementedFollowerServiceServer
	svc service.FollowService
}

func NewFollowerHandler(svc service.FollowService) *FollowerHandler {
	return &FollowerHandler{svc: svc}
}

// Helper to get UserID from API Gateway Metadata
func getUserID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no metadata provided")
	}
	// The Gateway passes the User ID in this header
	ids := md.Get("x-user-id")
	if len(ids) == 0 {
		return "", status.Error(codes.Unauthenticated, "user id not found in context")
	}
	return ids[0], nil
}

// 1. Follow
func (h *FollowerHandler) Follow(ctx context.Context, req *pb.FollowRequest) (*pb.FollowResponse, error) {
	followerID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	if followerID == req.TargetId {
		return nil, status.Error(codes.InvalidArgument, "cannot follow yourself")
	}

	if err := h.svc.FollowUser(ctx, followerID, req.TargetId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.FollowResponse{Status: "followed"}, nil
}

// 2. Unfollow
func (h *FollowerHandler) Unfollow(ctx context.Context, req *pb.UnfollowRequest) (*pb.UnfollowResponse, error) {
	followerID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.svc.UnfollowUser(ctx, followerID, req.TargetId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UnfollowResponse{Status: "unfollowed"}, nil
}

// 3. Recommendations
func (h *FollowerHandler) GetRecommendations(ctx context.Context, _ *emptypb.Empty) (*pb.RecommendationsResponse, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	recs, err := h.svc.GetRecommendations(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoRecs []*pb.RecommendationsResponse_RecUser
	for _, r := range recs {
		protoRecs = append(protoRecs, &pb.RecommendationsResponse_RecUser{
			UserId:            r.UserID,
			MutualConnections: r.MutualConnections,
		})
	}

	return &pb.RecommendationsResponse{Recommendations: protoRecs}, nil
}

// 4. Stats
func (h *FollowerHandler) GetStats(ctx context.Context, req *pb.StatsRequest) (*pb.StatsResponse, error) {
	// Note: req.UserId comes from the URL /{user_id}/stats
	stats, err := h.svc.GetUserStats(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.StatsResponse{
		UserId:         stats.UserId,
		FollowersCount: stats.FollowersCount,
		FollowingCount: stats.FollowingCount,
	}, nil
}

// 5. Followers List
func (h *FollowerHandler) GetFollowers(ctx context.Context, req *pb.GetListRequest) (*pb.GetListResponse, error) {
	ids, err := h.svc.GetFollowers(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetListResponse{UserIds: ids}, nil
}

// 6. Following List
func (h *FollowerHandler) GetFollowing(ctx context.Context, req *pb.GetListRequest) (*pb.GetListResponse, error) {
	ids, err := h.svc.GetFollowing(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetListResponse{UserIds: ids}, nil
}

// 7. Internal Check (IsFollowing)
func (h *FollowerHandler) IsFollowing(ctx context.Context, req *pb.IsFollowingRequest) (*pb.IsFollowingResponse, error) {
	isFollowing, err := h.svc.IsFollowing(ctx, req.FollowerId, req.FolloweeId)
	if err != nil {
		return nil, err
	}
	return &pb.IsFollowingResponse{IsFollowing: isFollowing}, nil
}
