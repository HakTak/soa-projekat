package model

import (
	"time"
)

// Comment predstavlja recenziju korisnika za neku turu
type Comment struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	TourID    string    `gorm:"column:tour_id" json:"tour_id"`
	UserID    string    `gorm:"column:user_id" json:"user_id"`
	Rating    int       `json:"rating"` // 1-5
	Content   string    `json:"content"`
	UserName  string    `json:"user_name"`
	TourDate  time.Time `json:"tour_date"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	ImageURL  string    `json:"image_url"`
}
