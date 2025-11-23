package service

import (
	"tour-service/internal/model"
	"tour-service/internal/repository"
)

type TourService struct {
	repo       *repository.TourRepository
	reviewRepo *repository.ReviewRepository
}

func NewTourService(repo *repository.TourRepository, reviewRepo *repository.ReviewRepository) *TourService {
	return &TourService{repo: repo, reviewRepo: reviewRepo}
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

func (s *TourService) UpdateTour(t *model.Tour) error {
	return s.repo.UpdateTour(t)
}

func (s *TourService) GetToursByUser(userId uint) ([]model.Tour, error) {
	return s.repo.GetToursByUser(userId)
}

// ****** Review service methods ******//
func (s *TourService) CreateReview(rev *model.Review) error {
	return s.reviewRepo.CreateReview(rev)
}

func (s *TourService) GetReviewsByTour(tourID string) ([]model.Review, error) {
	return s.reviewRepo.GetReviewsByTour(tourID)
}

func (s *TourService) DeleteReview(id string) error {
	return s.reviewRepo.DeleteReview(id)
}

func (s *TourService) GetAllReviews() ([]model.Review, error) {
	return s.reviewRepo.GetAllReviews()
}
