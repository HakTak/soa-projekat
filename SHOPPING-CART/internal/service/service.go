package service

import (
	"SHOPPING-CART/internal/model"
	"SHOPPING-CART/internal/repository"
	"context"
	"errors"
	"fmt"

	// Import the generated protobuf code for the Tour Service
	// Make sure this path matches your go.mod module name
	tourPb "SHOPPING-CART/common/genproto"
)

type ShoppingCartService struct {
	repo       *repository.CartRepository
	tourClient tourPb.TourServiceClient // gRPC Client interface
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
	// We need to fetch details (Name, Price) and check status
	tourResp, err := s.tourClient.GetTour(context.Background(), &tourPb.GetTourRequest{TourId: tourID})
	if err != nil {
		return nil, fmt.Errorf("failed to communicate with Tour service: %v", err)
	}

	// 2. VALIDATION: "Arhivirane ture se ne mogu kupiti"
	if tourResp.IsArchived {
		return nil, errors.New("this tour is archived or draft and cannot be purchased")
	}

	// 3. Get User's Cart
	cart, err := s.repo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	// 4. Check if item already exists (Prevent duplicates)
	for _, item := range cart.Items {
		if item.TourID == tourID {
			return nil, errors.New("tour is already in your cart")
		}
	}

	// 5. Create new Item
	newItem := model.OrderItem{
		CartID:   userID,
		TourID:   tourResp.Id,
		TourName: tourResp.Name,
		Price:    tourResp.Price,
	}

	// Add to list
	cart.Items = append(cart.Items, newItem)

	// 6. "Korpa računa ukupnu cenu..."
	cart.CalculateTotal()

	// 7. Save to DB
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

	// Filter items
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

	// 1. Delete the item from DB specifically
	if err := s.repo.RemoveItemByID(foundItem.ID); err != nil {
		return nil, err
	}

	// 2. Update struct locally to recalculate price
	cart.Items = keptItems
	cart.CalculateTotal()

	// 3. Update Cart Header (Total Price) in DB
	if err := s.repo.UpdateCart(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

// Checkout finalizes the purchase
func (s *ShoppingCartService) Checkout(userID string) ([]model.PurchaseToken, error) {
	// 1. Get Cart
	cart, err := s.repo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cannot checkout an empty cart")
	}

	// 2. Generate Tokens
	// "...za svaku stavku iz korpe dobija token"
	var tokens []model.PurchaseToken
	for _, item := range cart.Items {
		// UUIDs are usually generated in the Model's BeforeCreate hook,
		// but we create the struct here.
		token := model.PurchaseToken{
			UserID: userID,
			TourID: item.TourID,
			// Token ID and Secret are handled by GORM Hooks (model.go)
		}
		tokens = append(tokens, token)
	}

	// 3. Execute Transaction (Save Tokens + Clear Cart)
	if err := s.repo.ProcessCheckout(tokens, userID); err != nil {
		return nil, fmt.Errorf("checkout failed: %v", err)
	}

	return tokens, nil
}

func (s *ShoppingCartService) IsTourPurchased(tourId, touristId string) (bool, error) {
	return s.repo.IsTourPurchased(tourId, touristId)
}
