package repository

import (
	"tour-service/internal/model"

	"gorm.io/gorm"
)

type TourRepository struct {
	db *gorm.DB
}

func NewTourRepository(db *gorm.DB) *TourRepository {
	return &TourRepository{db}
}

func (r *TourRepository) CreateTour(t *model.Tour) error {
	return r.db.Create(t).Error
}

func (r *TourRepository) GetTour(id uint) (*model.Tour, error) {
	var tour model.Tour
	err := r.db.Preload("Keypoints").First(&tour, id).Error
	return &tour, err
}

func (r *TourRepository) GetAllTours() ([]model.Tour, error) {
	var tours []model.Tour
	err := r.db.Preload("Keypoints").Find(&tours).Error
	return tours, err
}

func (r *TourRepository) DeleteTour(id uint) error {
	return r.db.Delete(&model.Tour{}, id).Error
}
