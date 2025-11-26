package api

import (
	"context"
	"tour-service/internal/service"
	pb "tour-service/protobuf" // Import generated protobuf code

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TourGrpcHandler implements the generated TourServiceServer interface
type TourGrpcHandler struct {
	pb.UnimplementedTourServiceServer
	service *service.TourService
}

// NewTourGrpcHandler creates a new instance
func NewTourGrpcHandler(s *service.TourService) *TourGrpcHandler {
	return &TourGrpcHandler{service: s}
}

// GetTour is the function called by the Shopping Cart Service
func (h *TourGrpcHandler) GetTour(ctx context.Context, req *pb.GetTourRequest) (*pb.GetTourResponse, error) {
	// 1. Call your existing business logic to get the tour from DB
	tour, err := h.service.GetTour(req.TourId)
	if err != nil {
		// Return a gRPC-specific error code (NOT HTTP 404)
		return nil, status.Errorf(codes.NotFound, "tour with id %s not found", req.TourId)
	}

	// 2. CHECK STATUS LOGIC
	// Your requirement: "Archived tours cannot be purchased."
	// Usually, "DRAFT" tours shouldn't be bought either.

	isArchived := true // Default to unavailable

	// Only if the status is explicitly PUBLISHED, we allow purchase
	if tour.Status == "PUBLISHED" {
		isArchived = false
	}

	// 3. Map to Protobuf Response
	response := &pb.GetTourResponse{
		Id:         tour.ID,
		Name:       tour.Title, // Mapping 'Title' from DB to 'Name' in Proto
		Price:      tour.Price,
		IsArchived: isArchived, // This tells Shopping Cart if it's safe to buy
	}

	return response, nil
}
