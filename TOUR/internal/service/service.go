package service

import (
	"tour-service/internal/model"
	"tour-service/internal/repository"
)

type TourService struct {
	repo *repository.TourRepository
}

func NewTourService(repo *repository.TourRepository) *TourService {
	return &TourService{repo: repo}
}

func (s *TourService) CreateTour(t *model.Tour) error {
	return s.repo.CreateTour(t)
}

func (s *TourService) GetTour(id uint) (*model.Tour, error) {
	return s.repo.GetTour(id)
}

func (s *TourService) GetAllTours() ([]model.Tour, error) {
	return s.repo.GetAllTours()
}

func (s *TourService) DeleteTour(id uint) error {
	return s.repo.DeleteTour(id)
}
