package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tour struct {
	ID          string     `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Difficulty  string     `json:"difficulty"`
	Tags        string     `json:"tags"`
	Status      string     `json:"status"`
	Price       float64    `json:"price"`
	Keypoints   []Keypoint `gorm:"foreignKey:TourID" json:"keypoints" constraint:"OnDelete:CASCADE"`
}

func (t *Tour) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New().String()
	return
}
