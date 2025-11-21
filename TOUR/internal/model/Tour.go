package model

import (
	"time"

	"github.com/lib/pq"
)

type Tour struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	IDAuthor    string         `gorm:"column:id_author" json:"id_author"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	PublishedAt time.Time      `json:"published"`
	Dificulty   DificultyLevel `json:"dificulty_level"`
	Status      Status         `json:"status"`
	Tags        pq.StringArray `gorm:"type:text[]" json:"tags"`
}

type DificultyLevel string

const (
	DificultyEasy   DificultyLevel = "easy"
	DificultyMedium DificultyLevel = "medium"
	DificultyHard   DificultyLevel = "hard"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusDraft    Status = "draft"
	StatusArchived Status = "archived"
)

type Tag string

const (
	Cultural   Tag = "cultural"
	Food       Tag = "food"
	Attraction Tag = "attraction"
	Beach      Tag = "beach"
	Mountain   Tag = "mountain"
	Nature     Tag = "nature"
	Adventure  Tag = "adventure"
	Relax      Tag = "relax"
	Nightlife  Tag = "nightlife"
	Shopping   Tag = "shopping"
)
