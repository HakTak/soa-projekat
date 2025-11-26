package handlers

import (
	"context"

	// ✅ Update these imports to match your project structure
	"SHOPPING-CART/internal/model"
	"SHOPPING-CART/internal/service"

	// ✅ Import the generated code from COMMON
	pb "PROJEKAT/COMMON/shopping-cart/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ShoppingCartHandler struct {
	pb.UnimplementedShoppingCartServiceServer
	svc *service.ShoppingCartService
}

func NewShoppingCartHandler(svc *service.ShoppingCartService) *ShoppingCartHandler {
	return &ShoppingCartHandler{svc: svc}
}

// --- HELPER: Get User ID from Metadata (Same as Follower) ---
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

// --- HELPER: Map DB Model to Proto Message ---
// This keeps the main methods clean
func mapCartToProto(cart *model.ShoppingCart) *pb.CartResponse {
	var protoItems []*pb.CartItem

	for _, item := range cart.Items {
		protoItems = append(protoItems, &pb.CartItem{
			TourId:   item.TourID,
			TourName: item.TourName,
			Price:    item.Price,
		})
	}

	return &pb.CartResponse{
		UserId:     cart.UserID,
		Items:      protoItems,
		TotalPrice: cart.Total,
	}
}

// 1. GetCart
func (h *ShoppingCartHandler) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.CartResponse, error) {
	// Security: Always use the ID from the Token (Metadata), not just the URL
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	cart, err := h.svc.GetCart(userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return mapCartToProto(cart), nil
}

// 2. AddItem
func (h *ShoppingCartHandler) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.CartResponse, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	// Call Service
	cart, err := h.svc.AddItem(userID, req.TourId)
	if err != nil {
		// You might want to check specific errors here (e.g. "tour not found")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return mapCartToProto(cart), nil
}

// 3. RemoveItem
func (h *ShoppingCartHandler) RemoveItem(ctx context.Context, req *pb.RemoveItemRequest) (*pb.CartResponse, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	cart, err := h.svc.RemoveItem(userID, req.TourId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return mapCartToProto(cart), nil
}

// 4. Checkout
func (h *ShoppingCartHandler) Checkout(ctx context.Context, req *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	// Call Service
	tokens, err := h.svc.Checkout(userID)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	// Map DB Tokens to Proto Tokens
	var protoTokens []*pb.PurchaseToken
	for _, t := range tokens {
		protoTokens = append(protoTokens, &pb.PurchaseToken{
			TokenId: t.Token, // The secret token
			TourId:  t.TourID,
			UserId:  t.UserID,
		})
	}

	return &pb.CheckoutResponse{
		Success: true,
		Tokens:  protoTokens,
	}, nil
}
