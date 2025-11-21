package model

type Tour struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Difficulty  string     `json:"difficulty"`
	Tags        string     `json:"tags"`
	Status      string     `json:"status"`
	Price       float64    `json:"price"`
	Keypoints   []Keypoint `gorm:"foreignKey:TourID" json:"keypoints" constraint:"OnDelete:CASCADE"`
}
