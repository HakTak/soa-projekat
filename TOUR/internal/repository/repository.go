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

func (r *TourRepository) CreateTour(t *model.Tour) (*model.Tour, error) {
	if err := r.db.Create(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (r *TourRepository) GetTour(id string) (*model.Tour, error) {
	var tour model.Tour
	err := r.db.Preload("Keypoints").Preload("RouteOptions").First(&tour, "id = ?", id).Error
	return &tour, err
}

func (r *TourRepository) GetAllTours() ([]model.Tour, error) {
	var tours []model.Tour
	err := r.db.Preload("Keypoints").Preload("RouteOptions").Find(&tours).Error
	return tours, err
}

func (r *TourRepository) DeleteTour(id string) error {
	return r.db.Delete(&model.Tour{}, "id = ?", id).Error
}

func (r *TourRepository) UpdateTour(t *model.Tour) (*model.Tour, error) {
	if err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (r *TourRepository) GetToursByUser(userId string) ([]model.Tour, error) {
	var tours []model.Tour
	err := r.db.Preload("Keypoints").Preload("RouteOptions").Where("user_id = ?", userId).Find(&tours).Error
	return tours, err
}
