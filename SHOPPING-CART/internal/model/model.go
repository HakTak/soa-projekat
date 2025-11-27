package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ShoppingCart represents the user's active session cart
type ShoppingCart struct {
	UserID    string      `gorm:"primaryKey" json:"user_id"`
	Items     []OrderItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE" json:"items"`
	Total     float64     `json:"total"`
	UpdatedAt time.Time
}

// OrderItem represents a tour inside the cart
type OrderItem struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	CartID   string  `json:"-"` // Foreign Key to ShoppingCart
	TourID   string  `json:"tour_id"`
	TourName string  `json:"tour_name"`
	Price    float64 `json:"price"`
}

// PurchaseToken represents a purchased item (Receipt)
type PurchaseToken struct {
	ID       string    `gorm:"primaryKey" json:"id"`
	Token    string    `gorm:"uniqueIndex" json:"token"` // The actual token string
	UserID   string    `json:"user_id"`
	TourID   string    `json:"tour_id"`
	IssuedAt time.Time `json:"issued_at"`
	Status   string    `json:"status"`
}

// BeforeCreate generates a UUID for the token
func (pt *PurchaseToken) BeforeCreate(tx *gorm.DB) (err error) {
	pt.ID = uuid.New().String()
	pt.Token = uuid.New().String() // This is the "key" to unlock tour details
	pt.IssuedAt = time.Now()
	return
}

// Helper to recalculate price
func (sc *ShoppingCart) CalculateTotal() {
	var sum float64
	for _, item := range sc.Items {
		sum += item.Price
	}
	sc.Total = sum
}
