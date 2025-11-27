package service

import (
	"errors"
	"time"
	"tour-service/internal/model"
	"tour-service/internal/repository"
)

type TourService struct {
	repo              *repository.TourRepository
	reviewRepo        *repository.ReviewRepository
	tourExecutionRepo *repository.TourExecutionRepository
}

func NewTourService(repo *repository.TourRepository, reviewRepo *repository.ReviewRepository, tourExecutionRepo *repository.TourExecutionRepository) *TourService {
	return &TourService{repo: repo, reviewRepo: reviewRepo, tourExecutionRepo: tourExecutionRepo}
}

func (s *TourService) CreateTour(t *model.Tour) error {
	return s.repo.CreateTour(t)
}

func (s *TourService) GetTour(id string) (*model.Tour, error) {
	return s.repo.GetTour(id)
}

func (s *TourService) GetAllTours() ([]model.Tour, error) {
	return s.repo.GetAllTours()
}

func (s *TourService) DeleteTour(id string) error {
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

func (s *TourService) CreateTourExecution(teDto model.TourCreateExecutionDTO) (*model.TourExecution, error) {
	te := &model.TourExecution{
		TourId:          teDto.TourId,
		TouristId:       teDto.TouristId,
		Status:          "ACTIVE",
		StartedAt:       time.Now(),
		CurrentPosition: teDto.CurrentPosition,
	}

	tour, err1 := s.repo.GetTour(teDto.TourId)

	if err1 != nil {
		return nil, err1
	}

	if tour.Status == "DRAFT" {
		return nil, errors.New("Tura ne moze biti pokrenuta jer je u DRAFT statusu")
	}

	//ovde proveri da li je tura kupljena

	err := s.tourExecutionRepo.Create(te)
	return te, err
}

func (s *TourService) GetById(id string) (*model.TourExecution, error) {
	return s.tourExecutionRepo.GetById(id)
}

func (s *TourService) FinishOrAbandon(id, status string) (*model.TourExecution, error) {
	return s.tourExecutionRepo.FinishOrAbandon(id, model.TourExecutionStatus(status))
}

func (s *TourService) SetCurrentLocation(touristId string, curLoc model.TouristPosition) (int, error) {
	return s.tourExecutionRepo.SetCurrentLocation(touristId, curLoc)
}

func (s *TourService) CheckCloseKeyPoints(id string, tourId string) (*model.TourExecution, error) {
	kps, err1 := s.tourExecutionRepo.GetKeypointsByTour(tourId)

	if err1 != nil {
		return nil, err1
	}

	te, err2 := s.tourExecutionRepo.GetById(id)

	if err2 != nil {
		return nil, err2
	}

	kpId := model.FindNearbyKeypointId(kps, te.CurrentPosition.Latitude, te.CurrentPosition.Longitude)

	if kpId == "" {
		return nil, nil
	}

	completedKps := &model.CompletedKeypoint{
		KeyPointId:  kpId,
		CompletedAt: time.Now(),
	}

	te.CompledetKeyPonts = append(te.CompledetKeyPonts, *completedKps)

	updatedTe, err3 := s.tourExecutionRepo.CompleteKeyPoint(id, te.CompledetKeyPonts)

	if err3 != nil {
		return nil, err3
	}

	if len(updatedTe.CompledetKeyPonts) < len(kps) {
		return updatedTe, nil
	}

	return s.FinishOrAbandon(id, "COMPLETED")
}
