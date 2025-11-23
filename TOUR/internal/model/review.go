package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	TourID    string    `gorm:"type:uuid;index" json:"tour_id"`
	UserID    string    `gorm:"type:uuid;index" json:"user_id"`
	UserName  string    `json:"user_name"`
	Rating    int       `json:"rating"` // 1-5
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"` // kada je recenzija kreirana
	VisitedAt time.Time `json:"visited_at"` // kada je bio na turi
	ImageURL  string    `json:"image_url"`
}

func (r *Review) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New().String()
	return
}
