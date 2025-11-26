package handlers

import (
	"context"

	// ✅ Import imports
	"FOLLOWER/internal/service"
	pb "PROJEKAT/COMMON/follower/proto"
	"PROJEKAT/COMMON/utils" // <--- Import shared utils

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

// =================================================================
// SAFE HELPER: Retrieves claims from Context or Metadata (Fallback)
// =================================================================
func (h *FollowerHandler) getClaimsSafe(ctx context.Context) (map[string]interface{}, error) {
	// 1. Try getting claims from the shared utility
	claims := utils.ClaimsFromContext(ctx)

	// If utils returned a map, ensure it has the ID
	if claims != nil {
		if _, ok := claims["id"].(string); ok {
			return claims, nil
		}
	}

	// 2. FALLBACK: Read directly from gRPC Metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "No metadata provided")
	}

	newClaims := make(map[string]interface{})

	// Gateway sends "x-user-id" (canonical) or "user-id"
	// We check both variants to be safe
	if ids := md.Get("x-user-id"); len(ids) > 0 {
		newClaims["id"] = ids[0]
	} else if ids := md.Get("user-id"); len(ids) > 0 {
		newClaims["id"] = ids[0]
	}

	// Gateway sends "x-user-role"
	if roles := md.Get("x-user-role"); len(roles) > 0 {
		newClaims["role"] = roles[0]
	} else if roles := md.Get("user-role"); len(roles) > 0 {
		newClaims["role"] = roles[0]
	}

	// Final check: Do we have an ID?
	if _, ok := newClaims["id"]; !ok {
		return nil, status.Error(codes.Unauthenticated, "User ID not found in request")
	}

	return newClaims, nil
}

// =================================================================
// HANDLERS
// =================================================================

// 1. Follow
func (h *FollowerHandler) Follow(ctx context.Context, req *pb.FollowRequest) (*pb.FollowResponse, error) {
	// ✅ Use Safe Claims Extraction
	claims, err := h.getClaimsSafe(ctx)
	if err != nil {
		return nil, err
	}
	followerID := claims["id"].(string) // Safe because getClaimsSafe verified it

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
	// ✅ Use Safe Claims Extraction
	claims, err := h.getClaimsSafe(ctx)
	if err != nil {
		return nil, err
	}
	followerID := claims["id"].(string)

	if err := h.svc.UnfollowUser(ctx, followerID, req.TargetId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UnfollowResponse{Status: "unfollowed"}, nil
}

// 3. Recommendations
func (h *FollowerHandler) GetRecommendations(ctx context.Context, _ *emptypb.Empty) (*pb.RecommendationsResponse, error) {
	// ✅ Use Safe Claims Extraction
	claims, err := h.getClaimsSafe(ctx)
	if err != nil {
		return nil, err
	}
	userID := claims["id"].(string)

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
	// Note: We use req.UserId from URL here (viewing someone else's stats)
	// If you wanted to restrict this to only allow viewing your own stats,
	// you would call h.getClaimsSafe(ctx) here too.
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
