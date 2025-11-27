package handlers

import (
	"context"

	// ✅ Ensure imports are correct
	pb "PROJEKAT/COMMON/shopping-cart/proto"
	"PROJEKAT/COMMON/utils" // Import shared utils
	"SHOPPING-CART/internal/model"
	"SHOPPING-CART/internal/service"

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

// =================================================================
// SAFE HELPER: Retrieves claims from Context or Metadata (Fallback)
// =================================================================
func (h *ShoppingCartHandler) getClaimsSafe(ctx context.Context) (map[string]interface{}, error) {
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
	if ids := md.Get("x-user-id"); len(ids) > 0 {
		newClaims["id"] = ids[0]
	} else if ids := md.Get("user-id"); len(ids) > 0 {
		newClaims["id"] = ids[0]
	}

	// Final check: Do we have an ID?
	if _, ok := newClaims["id"]; !ok {
		return nil, status.Error(codes.Unauthenticated, "alo required (User ID missing)")
	}

	return newClaims, nil
}

// --- HELPER: Map DB Model to Proto Message ---
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

// =================================================================
// HANDLERS
// =================================================================

// 1. GetCart
func (h *ShoppingCartHandler) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.CartResponse, error) {
	// ✅ USE SAFE HELPER
	claims, err := h.getClaimsSafe(ctx)
	if err != nil {
		return nil, err
	}
	userID := claims["id"].(string)

	cart, err := h.svc.GetCart(userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return mapCartToProto(cart), nil
}

// 2. AddItem
func (h *ShoppingCartHandler) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.CartResponse, error) {
	// ✅ USE SAFE HELPER
	claims, err := h.getClaimsSafe(ctx)
	if err != nil {
		return nil, err
	}
	userID := claims["id"].(string)

	// Call Service
	cart, err := h.svc.AddItem(userID, req.TourId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return mapCartToProto(cart), nil
}

// 3. RemoveItem
func (h *ShoppingCartHandler) RemoveItem(ctx context.Context, req *pb.RemoveItemRequest) (*pb.CartResponse, error) {
	// ✅ USE SAFE HELPER
	claims, err := h.getClaimsSafe(ctx)
	if err != nil {
		return nil, err
	}
	userID := claims["id"].(string)

	cart, err := h.svc.RemoveItem(userID, req.TourId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return mapCartToProto(cart), nil
}

// 4. Checkout
func (h *ShoppingCartHandler) Checkout(ctx context.Context, req *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {
	// ✅ USE SAFE HELPER
	claims, err := h.getClaimsSafe(ctx)
	if err != nil {
		return nil, err
	}
	userID := claims["id"].(string)

	// Call Service
	tokens, err := h.svc.Checkout(userID)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	// Map DB Tokens to Proto Tokens
	var protoTokens []*pb.PurchaseToken
	for _, t := range tokens {
		protoTokens = append(protoTokens, &pb.PurchaseToken{
			TokenId: t.Token,
			TourId:  t.TourID,
			UserId:  t.UserID,
		})
	}

	return &pb.CheckoutResponse{
		Success: true,
		Tokens:  protoTokens,
	}, nil
}

// Add this method to ShoppingCartHandler
func (h *ShoppingCartHandler) GetPurchasedTours(ctx context.Context, req *pb.GetPurchasedToursRequest) (*pb.GetPurchasedToursResponse, error) {
	// 1. Safe Authentication Check
	claims, err := h.getClaimsSafe(ctx)
	if err != nil {
		return nil, err
	}
	userID := claims["id"].(string)

	// 2. Call Service
	tokens, err := h.svc.GetPurchasedTours(userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 3. Map DB Models to Proto Messages
	var protoTokens []*pb.PurchaseToken
	for _, t := range tokens {
		protoTokens = append(protoTokens, &pb.PurchaseToken{
			TokenId: t.Token, // The secret token
			TourId:  t.TourID,
			UserId:  t.UserID,
			Status:  t.Status,
		})
	}

	return &pb.GetPurchasedToursResponse{
		Tokens: protoTokens,
	}, nil
}
