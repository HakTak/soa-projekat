package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tour struct {
	ID           string        `gorm:"type:uuid;primaryKey" json:"id"`
	UserName     string        `json:"user_name"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Difficulty   string        `json:"difficulty"`
	Tags         string        `json:"tags"`
	Status       string        `json:"status"`
	Price        float64       `json:"price"`
	Distance     float64       `json:"distance"` // duzina ture u kilometrima
	Duration     time.Duration `json:"duration"` // trajanje ture u milisekundama u bazi bar
	CreatedAt    time.Time     `json:"created_at"`
	PublishedAt  time.Time     `json:"published_at"`
	ArchivedAt   time.Time     `json:"archived_at"`
	Keypoints    []Keypoint    `gorm:"foreignKey:TourID;constraint:OnDelete:CASCADE" json:"keypoints"`
	RouteOptions []RouteOption `gorm:"foreignKey:TourID;constraint:OnDelete:CASCADE" json:"route_options"`
	Sales        int64         `json:"sales" gorm:"default:0"`
}

func (t *Tour) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New().String()
	return
}
