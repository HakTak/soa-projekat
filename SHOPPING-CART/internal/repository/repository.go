package repository

import (
	"SHOPPING-CART/internal/model"
	"errors"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

// GetOrCreateCart fetches the cart with items. If it doesn't exist, it creates an empty one.
func (r *CartRepository) GetOrCreateCart(userID string) (*model.ShoppingCart, error) {
	var cart model.ShoppingCart

	// Preload("Items") fetches the OrderItems associated with this Cart
	err := r.db.Preload("Items").First(&cart, "user_id = ?", userID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new cart if not found
		cart = model.ShoppingCart{
			UserID: userID,
			Items:  []model.OrderItem{},
			Total:  0,
		}
		if err := r.db.Create(&cart).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &cart, nil
}

// UpdateCart saves changes to the cart (updated total) and its items (added items)
func (r *CartRepository) UpdateCart(cart *model.ShoppingCart) error {
	// Full save of the cart and its associations
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(cart).Error
}

// RemoveItemByID deletes a specific item from the database
func (r *CartRepository) RemoveItemByID(itemID uint) error {
	return r.db.Delete(&model.OrderItem{}, itemID).Error
}

// ProcessCheckout performs the "Purchase": Saves tokens and clears cart items in ONE transaction.
func (r *CartRepository) ProcessCheckout(tokens []model.PurchaseToken, userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Save all generated tokens (receipts)
		if err := tx.Create(&tokens).Error; err != nil {
			return err
		}

		// 2. Clear all items from the user's cart
		// We delete items where CartID == userID
		if err := tx.Where("cart_id = ?", userID).Delete(&model.OrderItem{}).Error; err != nil {
			return err
		}

		// 3. Reset the Cart's TotalPrice to 0
		if err := tx.Model(&model.ShoppingCart{}).Where("user_id = ?", userID).Update("total", 0).Error; err != nil {
			return err
		}

		return nil // Commit the transaction
	})
}

// Add this method to CartRepository
func (r *CartRepository) GetPurchasedTokens(userID string) ([]model.PurchaseToken, error) {
	var tokens []model.PurchaseToken

	// Fetch all tokens belonging to this user
	result := r.db.Where("user_id = ?", userID).Find(&tokens)

	if result.Error != nil {
		return nil, result.Error
	}

	return tokens, nil
}
