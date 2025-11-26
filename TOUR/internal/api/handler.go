package api

import (
	pb "PROJEKAT/COMMON/tour/proto"
	"PROJEKAT/COMMON/utils"
	"context"
	"time"
	"tour-service/internal/model"
	"tour-service/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TourGRPCServer struct {
	pb.UnimplementedTourServiceServer
	svc *service.TourService
}

func NewTourGRPCServer(svc *service.TourService) *TourGRPCServer {
	return &TourGRPCServer{svc: svc}
}

// ========================
// Tour Methods
// ========================

func (s *TourGRPCServer) CreateTour(ctx context.Context, req *pb.CreateTourRequest) (*pb.TourResponse, error) {
	if err := utils.Authorize(ctx, "GUIDE"); err != nil {
		return nil, status.Error(codes.PermissionDenied, "Access Denied")
	}
	t := &model.Tour{
		UserName:    req.Tour.UserName,
		Title:       req.Tour.Title,
		Description: req.Tour.Description,
		Difficulty:  req.Tour.Difficulty,
		Tags:        req.Tour.Tags,
		Status:      req.Tour.Status.String(),
		Price:       req.Tour.Price,
		Distance:    req.Tour.Distance,
		Duration:    time.Duration(req.Tour.Duration) * time.Millisecond,
	}

	for _, kp := range req.Tour.Keypoints {
		t.Keypoints = append(t.Keypoints, model.Keypoint{
			Title:       kp.Title,
			Latitude:    kp.Latitude,
			Longitude:   kp.Longitude,
			Description: kp.Description,
			ImageURL:    kp.ImageUrl,
		})
	}

	for _, ro := range req.Tour.RouteOptions {
		t.RouteOptions = append(t.RouteOptions, model.RouteOption{
			Mode:     model.RouteMode(ro.Mode.String()),
			Duration: time.Duration(ro.Duration) * time.Millisecond,
		})
	}

	createdTour, err := s.svc.CreateTour(t)
	if err != nil {
		return nil, err
	}

	return mapTourToResponse(createdTour), nil
}

func (s *TourGRPCServer) UpdateTour(ctx context.Context, req *pb.UpdateTourRequest) (*pb.TourResponse, error) {
	if err := utils.Authorize(ctx, "GUIDE"); err != nil {
		return nil, status.Error(codes.PermissionDenied, "Access Denied")
	}
	t := &model.Tour{
		ID:          req.Tour.Id,
		UserName:    req.Tour.UserName,
		Title:       req.Tour.Title,
		Description: req.Tour.Description,
		Difficulty:  req.Tour.Difficulty,
		Tags:        req.Tour.Tags,
		Status:      req.Tour.Status.String(),
		Price:       req.Tour.Price,
		Distance:    req.Tour.Distance,
		Duration:    time.Duration(req.Tour.Duration) * time.Millisecond,
	}

	for _, kp := range req.Tour.Keypoints {
		t.Keypoints = append(t.Keypoints, model.Keypoint{
			ID:          kp.Id,
			Title:       kp.Title,
			Latitude:    kp.Latitude,
			Longitude:   kp.Longitude,
			Description: kp.Description,
			ImageURL:    kp.ImageUrl,
		})
	}

	for _, ro := range req.Tour.RouteOptions {
		t.RouteOptions = append(t.RouteOptions, model.RouteOption{
			ID:       ro.Id,
			Mode:     model.RouteMode(ro.Mode.String()),
			Duration: time.Duration(ro.Duration) * time.Millisecond,
		})
	}

	updatedTour, err := s.svc.UpdateTour(t)
	if err != nil {
		return nil, err
	}

	return mapTourToResponse(updatedTour), nil
}

func (s *TourGRPCServer) GetTour(ctx context.Context, req *pb.GetTourRequest) (*pb.TourResponse, error) {
	t, err := s.svc.GetTour(req.Id)
	if err != nil {
		return nil, err
	}
	return mapTourToResponse(t), nil
}

func (s *TourGRPCServer) GetAllTours(ctx context.Context, _ *emptypb.Empty) (*pb.GetAllToursResponse, error) {
	tours, err := s.svc.GetAllTours()
	if err != nil {
		return nil, err
	}
	return mapAllToursToResponse(tours), nil
}

func (s *TourGRPCServer) GetToursByUser(ctx context.Context, req *pb.GetToursByUserRequest) (*pb.GetAllToursResponse, error) {
	tours, err := s.svc.GetToursByUser(req.UserName)
	if err != nil {
		return nil, err
	}
	return mapAllToursToResponse(tours), nil
}

func (s *TourGRPCServer) DeleteTour(ctx context.Context, req *pb.DeleteTourRequest) (*pb.TourResponse, error) {
	if err := utils.Authorize(ctx, "GUIDE"); err != nil {
		return nil, status.Error(codes.PermissionDenied, "Access Denied")
	}
	t, err := s.svc.GetTour(req.Id)
	if err != nil {
		return nil, err
	}

	if err := s.svc.DeleteTour(req.Id); err != nil {
		return nil, err
	}

	return mapTourToResponse(t), nil
}

// ========================
// Review Methods
// ========================

func (s *TourGRPCServer) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.ReviewResponse, error) {
	claims := utils.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, status.Error(codes.Unauthenticated, "User is not authenticated")
	}
	r := &model.Review{
		TourID:    req.Review.TourId,
		UserName:  req.Review.UserName,
		Rating:    int(req.Review.Rating),
		Comment:   req.Review.Comment,
		VisitedAt: req.Review.VisitedAt.AsTime(),
		ImageURL:  req.Review.ImageUrl,
	}

	if err := s.svc.CreateReview(r); err != nil {
		return nil, err
	}

	return mapReviewToResponse(r), nil
}

func (s *TourGRPCServer) GetReviewsByTour(ctx context.Context, req *pb.GetReviewsByTourRequest) (*pb.GetAllReviewsResponse, error) {
	revs, err := s.svc.GetReviewsByTour(req.TourId)
	if err != nil {
		return nil, err
	}

	return mapAllReviewsToResponse(revs), nil
}

func (s *TourGRPCServer) GetAllReviews(ctx context.Context, _ *emptypb.Empty) (*pb.GetAllReviewsResponse, error) {
	revs, err := s.svc.GetAllReviews()
	if err != nil {
		return nil, err
	}

	return mapAllReviewsToResponse(revs), nil
}

func (s *TourGRPCServer) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.ReviewResponse, error) {
	claims := utils.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, status.Error(codes.Unauthenticated, "User is not authenticated")
	}
	revList, err := s.svc.GetReviewsByTour(req.Id) // Implement GetReviewByID if needed
	if err != nil || len(revList) == 0 {
		return nil, err
	}

	if err := s.svc.DeleteReview(req.Id); err != nil {
		return nil, err
	}

	return mapReviewToResponse(&revList[0]), nil
}

// ========================
// Mapping Helpers
// ========================

func mapTourToResponse(t *model.Tour) *pb.TourResponse {
	return &pb.TourResponse{
		Tour: &pb.Tour{
			Id:           t.ID,
			UserName:     t.UserName,
			Title:        t.Title,
			Description:  t.Description,
			Difficulty:   t.Difficulty,
			Tags:         t.Tags,
			Status:       mapTourStatusToProto(t.Status),
			Price:        t.Price,
			Distance:     t.Distance,
			Duration:     int64(t.Duration / time.Millisecond),
			CreatedAt:    timestamppb.New(t.CreatedAt),
			PublishedAt:  timestamppb.New(t.PublishedAt),
			ArchivedAt:   timestamppb.New(t.ArchivedAt),
			Keypoints:    mapKeypointsToProto(t.Keypoints),
			RouteOptions: mapRouteOptionsToProto(t.RouteOptions),
		},
	}
}

func mapAllToursToResponse(tours []model.Tour) *pb.GetAllToursResponse {
	resp := &pb.GetAllToursResponse{}
	for _, t := range tours {
		resp.Tours = append(resp.Tours, mapTourToResponse(&t).Tour)
	}
	return resp
}

func mapReviewToResponse(r *model.Review) *pb.ReviewResponse {
	return &pb.ReviewResponse{
		Review: &pb.Review{
			Id:        r.ID,
			TourId:    r.TourID,
			UserName:  r.UserName,
			Rating:    int32(r.Rating),
			Comment:   r.Comment,
			CreatedAt: timestamppb.New(r.CreatedAt),
			VisitedAt: timestamppb.New(r.VisitedAt),
			ImageUrl:  r.ImageURL,
		},
	}
}

func mapAllReviewsToResponse(revs []model.Review) *pb.GetAllReviewsResponse {
	resp := &pb.GetAllReviewsResponse{}
	for _, r := range revs {
		resp.Reviews = append(resp.Reviews, mapReviewToResponse(&r).Review)
	}
	return resp
}

func mapKeypointsToProto(kps []model.Keypoint) []*pb.Keypoint {
	res := []*pb.Keypoint{}
	for _, kp := range kps {
		res = append(res, &pb.Keypoint{
			Id:          kp.ID,
			TourId:      kp.TourID,
			Title:       kp.Title,
			Latitude:    kp.Latitude,
			Longitude:   kp.Longitude,
			Description: kp.Description,
			ImageUrl:    kp.ImageURL,
		})
	}
	return res
}

func mapRouteOptionsToProto(ros []model.RouteOption) []*pb.RouteOption {
	res := []*pb.RouteOption{}
	for _, ro := range ros {
		res = append(res, &pb.RouteOption{
			Id:       ro.ID,
			TourId:   ro.TourID,
			Mode:     mapRouteModeToProto(ro.Mode),
			Duration: int64(ro.Duration / time.Millisecond),
		})
	}
	return res
}

func mapRouteModeToProto(mode model.RouteMode) pb.RouteMode {
	switch mode {
	case model.Walking:
		return pb.RouteMode_WALKING
	case model.Cycling:
		return pb.RouteMode_CYCLING
	case model.Driving:
		return pb.RouteMode_DRIVING
	default:
		return pb.RouteMode_WALKING
	}
}

func mapTourStatusToProto(status string) pb.TourStatus {
	switch status {
	case "PUBLISHED":
		return pb.TourStatus_PUBLISHED
	case "ARCHIVED":
		return pb.TourStatus_ARCHIVED
	case "DRAFT":
		return pb.TourStatus_DRAFT
	default:
		return pb.TourStatus_DRAFT
	}
}
