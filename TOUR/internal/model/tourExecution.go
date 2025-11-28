package model

import (
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TourExecutionStatus string

const (
	StatusActive    TourExecutionStatus = "ACTIVE"
	StatusCompleted TourExecutionStatus = "COMPLETED"
	StatusAbandoned TourExecutionStatus = "ABANDONED"
)

type CompletedKeypoint struct {
	KeyPointId  string    `json:"key_point_id"`
	CompletedAt time.Time `json:"completed_at"`
}

type TouristPosition struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type TourExecution struct {
	Id                uuid.UUID           `json:"id" gorm:"type:uuid;primaryKey"`
	TourId            string              `json:"tour_id" gorm:"not null;unique"`
	TouristId         string              `json:"tourist_id" gorm:"not null;unique"`
	Status            TourExecutionStatus `json:"status"`
	StartedAt         time.Time           `json:"started_at"`
	LastActivityAt    time.Time           `json:"last_activity_at"`
	FinishedAt        time.Time           `json:"finished_at"`
	CompledetKeyPonts []CompletedKeypoint `json:"completed_key_points"`
	CurrentPosition   TouristPosition     `json:"current_position"`
}

func (u *TourExecution) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.New()
	return
}

type TourCreateExecutionDTO struct {
	TourId          string          `json:"tour_id" gorm:"not null;unique"`
	TouristId       string          `json:"tourist_id" gorm:"not null;unique"`
	CurrentPosition TouristPosition `json:"current_position"`
}

func FindNearbyKeypointId(keypoints []Keypoint, touristLat, touristLon float64) string {
	const hitRadiusMeters = 50.0

	for _, kp := range keypoints {
		distance := calculateHaversineDistance(touristLat, touristLon, kp.Latitude, kp.Longitude)

		if distance <= hitRadiusMeters {
			return kp.ID
		}
	}

	return ""
}

func calculateHaversineDistance(touristLat, touristLon, destLat, destLon float64) float64 {
	const R = 6371000

	phiT := touristLat * math.Pi / 180
	phiD := destLat * math.Pi / 180
	deltaPhi := (destLat - touristLat) * math.Pi / 180
	deltaLambda := (destLon - touristLon) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phiT)*math.Cos(phiD)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return distance
}
