package service

import (
	"SHOPPING-CART/internal/model"
	"SHOPPING-CART/internal/repository"
	"context"
	"errors"
	"fmt"

	// ✅ USE THE COMMON IMPORT (Matches your proto option go_package)
	tourPb "PROJEKAT/COMMON/tour/proto"
)

type ShoppingCartService struct {
	repo       *repository.CartRepository
	tourClient tourPb.TourServiceClient
}

func NewShoppingCartService(repo *repository.CartRepository, tourClient tourPb.TourServiceClient) *ShoppingCartService {
	return &ShoppingCartService{
		repo:       repo,
		tourClient: tourClient,
	}
}

// GetCart returns the cart for a user
func (s *ShoppingCartService) GetCart(userID string) (*model.ShoppingCart, error) {
	return s.repo.GetOrCreateCart(userID)
}

// AddItem adds a tour to the cart
func (s *ShoppingCartService) AddItem(userID, tourID string) (*model.ShoppingCart, error) {
	// 1. CALL TOUR SERVICE (gRPC)
	// Note: The proto field is now 'id', not 'tour_id'
	tourResp, err := s.tourClient.GetTour(context.Background(), &tourPb.GetTourRequest{Id: tourID})
	if err != nil {
		return nil, fmt.Errorf("failed to communicate with Tour service: %v", err)
	}

	// Safety check if Tour is nil inside the response
	if tourResp.Tour == nil {
		return nil, errors.New("tour not found")
	}

	// 2. VALIDATION: Check Enum Status
	// "Arhivirane ture se ne mogu kupiti" (Also probably not DRAFTs)
	if tourResp.Tour.Status == tourPb.TourStatus_ARCHIVED || tourResp.Tour.Status == tourPb.TourStatus_DRAFT {
		return nil, errors.New("this tour is archived or draft and cannot be purchased")
	}

	// 3. Get User's Cart
	cart, err := s.repo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	// 4. Check if item already exists
	for _, item := range cart.Items {
		if item.TourID == tourID {
			return nil, errors.New("tour is already in your cart")
		}
	}

	// 5. Create new Item
	// Note: We access fields via tourResp.Tour.Title / Price
	newItem := model.OrderItem{
		CartID:   userID,
		TourID:   tourResp.Tour.Id,
		TourName: tourResp.Tour.Title, // Proto field is 'title', mapped to 'TourName'
		Price:    tourResp.Tour.Price,
	}

	// Add to list
	cart.Items = append(cart.Items, newItem)

	// 6. Recalculate
	cart.CalculateTotal()

	// 7. Save
	if err := s.repo.UpdateCart(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// RemoveItem removes a tour from the cart
func (s *ShoppingCartService) RemoveItem(userID, tourID string) (*model.ShoppingCart, error) {
	cart, err := s.repo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	var keptItems []model.OrderItem
	var foundItem *model.OrderItem

	for i := range cart.Items {
		if cart.Items[i].TourID == tourID {
			foundItem = &cart.Items[i]
		} else {
			keptItems = append(keptItems, cart.Items[i])
		}
	}

	if foundItem == nil {
		return nil, errors.New("item not found in cart")
	}

	if err := s.repo.RemoveItemByID(foundItem.ID); err != nil {
		return nil, err
	}

	cart.Items = keptItems
	cart.CalculateTotal()

	if err := s.repo.UpdateCart(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// Add this method to ShoppingCartService
func (s *ShoppingCartService) GetPurchasedTours(userID string) ([]model.PurchaseToken, error) {
	return s.repo.GetPurchasedTokens(userID)
}

func (s *ShoppingCartService) Checkout(userID string) ([]model.PurchaseToken, error) {
	// 1. Get Cart
	cart, err := s.repo.GetOrCreateCart(userID)
	if err != nil || len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// =========================================================
	// SAGA STEP 1: PREPARE (Local Transaction)
	// Create tokens but mark them as PENDING.
	// If system crashes here, user has no valid tokens. Safe.
	// =========================================================
	var tokens []model.PurchaseToken
	for _, item := range cart.Items {
		tokens = append(tokens, model.PurchaseToken{
			UserID: userID,
			TourID: item.TourID,
			Status: "PENDING",
		})
	}

	if err := s.repo.CreatePendingTokens(tokens); err != nil {
		return nil, fmt.Errorf("checkout init failed: %v", err)
	}

	// =========================================================
	// SAGA STEP 2: ACTION (Remote Transaction)
	// Call Tour Service to increment sales counter.
	// =========================================================
	for _, item := range cart.Items {
		_, err := s.tourClient.UpdateTourSales(context.Background(), &tourPb.UpdateTourSalesRequest{
			TourId: item.TourID,
		})

		if err != nil {
			// !!! ERROR DETECTED !!!
			// The Tour Service said NO (or is down).

			// =========================================================
			// SAGA COMPENSATION (Rollback)
			// Undo Step 1 by deleting the pending tokens.
			// =========================================================
			fmt.Printf("Saga failure for tour %s: %v. Rolling back...\n", item.TourID, err)
			_ = s.repo.AbortPurchase(userID)

			return nil, fmt.Errorf("purchase failed: could not confirm tour availability. Cart has not been charged.")
		}
	}

	// =========================================================
	// SAGA STEP 3: COMMIT (Local Transaction)
	// Everything went well. Make tokens valid and clear cart.
	// =========================================================
	if err := s.repo.ConfirmPurchase(userID); err != nil {
		// This is a critical/rare edge case (DB died right at the end)
		return nil, fmt.Errorf("system error during finalization: %v", err)
	}

	// Return the confirmed tokens
	for i := range tokens {
		tokens[i].Status = "CONFIRMED"
	}
	return tokens, nil
}
