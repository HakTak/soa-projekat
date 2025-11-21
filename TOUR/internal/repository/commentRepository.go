package repository

import (
	"context"
	"errors"
	"tour/internal/model"

	"gorm.io/gorm"
)

var ErrCommentNotFound = errors.New("comment not found")

type CommentRepository interface {
	Create(ctx context.Context, c *model.Comment) error
	GetByTourID(ctx context.Context, tourID string) ([]*model.Comment, error)
}

type gormCommentRepo struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &gormCommentRepo{db: db}
}

// Create novi komentar
func (r *gormCommentRepo) Create(ctx context.Context, c *model.Comment) error {
	return r.db.WithContext(ctx).Create(c).Error
}

// GetByTourID vraća sve komentare za jednu turu
func (r *gormCommentRepo) GetByTourID(ctx context.Context, tourID string) ([]*model.Comment, error) {
	var comments []*model.Comment
	if err := r.db.WithContext(ctx).Where("tour_id = ?", tourID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
