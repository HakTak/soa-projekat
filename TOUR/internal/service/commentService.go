package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tour/internal/model"
	"tour/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentService struct {
	repo     repository.CommentRepository
	tourRepo repository.TourRepository
}

func NewCommentService(repo repository.CommentRepository, tourRepo repository.TourRepository) *CommentService {
	return &CommentService{repo: repo, tourRepo: tourRepo}
}

// Create novi komentar
func (s *CommentService) Create(ctx context.Context, tourID, userID, userName, content, imageURL string, rating int, tourDate time.Time) (*model.Comment, error) {

	_, err := s.tourRepo.GetByID(ctx, tourID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("tour with id %s does not exist", tourID)
		}
		return nil, err
	}

	comment := &model.Comment{
		ID:        uuid.NewString(),
		TourID:    tourID,
		UserID:    userID,
		UserName:  userName,
		Content:   content,
		ImageURL:  imageURL,
		Rating:    rating,
		TourDate:  tourDate,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

// GetByTourID vraća sve komentare za turu
func (s *CommentService) GetByTourID(ctx context.Context, tourID string) ([]*model.Comment, error) {
	return s.repo.GetByTourID(ctx, tourID)
}
