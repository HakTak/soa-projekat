package repository

import (
	"time"
	"tour-service/internal/model"

	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TourExecutionRepository struct {
	db *gorm.DB
}

func NewTourExecutionRepository(db *gorm.DB) *TourExecutionRepository {
	return &TourExecutionRepository{db: db}
}

func (r *TourExecutionRepository) Create(te *model.TourExecution) error {
	return r.db.Create(te).Error
}

func (r *TourExecutionRepository) GetById(id string) (*model.TourExecution, error) {
	var te model.TourExecution
	err := r.db.First(&te, "id = ?", id).Error
	return &te, err
}

func (r *TourExecutionRepository) GetByTourUser(tourId, touristId string) (*model.TourExecution, error) {
	var te model.TourExecution
	err := r.db.First(&te, "tour_id = ? and tourist_id", tourId, touristId).Error
	return &te, err
}

func (r *TourExecutionRepository) FinishOrAbandon(id string, status model.TourExecutionStatus) (*model.TourExecution, error) {
	uuidId := uuid.MustParse(id)
	res := r.db.Save(&model.TourExecution{Id: uuidId, Status: status, LastActivityAt: time.Now(), FinishedAt: time.Now()})

	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, errors.New(("user not found"))
	}

	return r.GetById(id)
}

func (r *TourExecutionRepository) GetKeypointsByTour(tourId string) ([]model.Keypoint, error) {
	var kps []model.Keypoint
	err := r.db.Find(&kps, "tourId = ?", tourId).Error
	return kps, err
}

func (r *TourExecutionRepository) SetCurrentLocation(touristId string, curLoc model.TouristPosition) (int, error) {
	res := r.db.Model(&model.TourExecution{}).
		Where("tourist_id = ? and status = ?", touristId, "ACTIVE").
		Update("current_position", curLoc)

	if res.Error != nil {
		return 0, res.Error
	}

	if res.RowsAffected == 0 {
		return 0, errors.New(("user not found"))
	}

	return int(res.RowsAffected), nil
}

func (r *TourExecutionRepository) GetCurrentLocation(id string) (*model.TouristPosition, error) {
	res := r.db.Model(&model.TourExecution{}).
		Where("id = ?", id).
		Update("last_activity_at", time.Now())

	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, errors.New(("user not found"))
	}

	te, err := r.GetById(id)

	return &te.CurrentPosition, err

}

func (r *TourExecutionRepository) CompleteKeyPoint(id string, updatedCKP []model.CompletedKeypoint) (*model.TourExecution, error) {
	res := r.db.Model(&model.TourExecution{}).
		Where("id = ?", id).
		Update("completed_key_points", updatedCKP)

	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, errors.New(("user not found"))
	}

	return r.GetById(id)
}
