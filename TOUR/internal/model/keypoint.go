package model

type Keypoint struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	TourID      uint    `json:"tourId"`
	Title       string  `json:"title"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
	ImageURL    string  `json:"imageUrl"`
}
