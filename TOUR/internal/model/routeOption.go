package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RouteOption struct {
	ID       string        `gorm:"type:uuid;primaryKey" json:"id"`
	TourID   string        `gorm:"type:uuid" json:"tourId"`
	Mode     RouteMode     `json:"mode"`
	Duration time.Duration `json:"duration"`
}

type RouteMode string

const (
	Walking RouteMode = "WALKING"
	Cycling RouteMode = "CYCLING"
	Driving RouteMode = "DRIVING"
)

func (r *RouteOption) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New().String()
	return
}
