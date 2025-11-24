package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Keypoint struct {
	ID          string  `gorm:"type:uuid;primaryKey" json:"id"`
	TourID      string  `gorm:"type:uuid" json:"tourId"`
	Title       string  `json:"title"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
	ImageURL    string  `json:"imageUrl"`
}

func (k *Keypoint) BeforeCreate(tx *gorm.DB) (err error) {
	k.ID = uuid.New().String()
	return
}
