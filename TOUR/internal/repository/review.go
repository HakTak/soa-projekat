package repository

import (
	"tour-service/internal/model"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	db *gorm.DB
}

func CreateReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db}
}

func (r *ReviewRepository) CreateReview(rev *model.Review) error {
	return r.db.Create(rev).Error
}

func (r *ReviewRepository) GetReviewsByTour(tourID string) ([]model.Review, error) {
	var reviews []model.Review
	err := r.db.Where("tour_id = ?", tourID).Find(&reviews).Error
	return reviews, err
}

func (r *ReviewRepository) DeleteReview(id string) error {
	return r.db.Delete(&model.Review{}, "id = ?", id).Error
}
