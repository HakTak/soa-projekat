package repository

import (
	"context"
	"errors"
	"tour/internal/model"

	"gorm.io/gorm"
)

var (
	ErrTourNotFound      = errors.New("tour not found")
	ErrTourAlreadyExists = errors.New("tour already exists")
)

type TourRepository interface {
	Create(ctx context.Context, t *model.Tour) error
	GetByID(ctx context.Context, id string) (*model.Tour, error)
	GetByAuthor(ctx context.Context, authorID string) ([]*model.Tour, error)
	Update(ctx context.Context, t *model.Tour) error
	Delete(ctx context.Context, id string) error
}

type gormTourRepo struct {
	db *gorm.DB
}

func NewGormTourRepo(db *gorm.DB) TourRepository {
	return &gormTourRepo{db: db}
}

func (r *gormTourRepo) Create(ctx context.Context, t *model.Tour) error {
	if err := r.db.WithContext(ctx).Create(t).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrTourAlreadyExists
		}
		return err
	}
	return nil
}

func (r *gormTourRepo) GetByID(ctx context.Context, id string) (*model.Tour, error) {
	var t model.Tour
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTourNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *gormTourRepo) GetByAuthor(ctx context.Context, authorID string) ([]*model.Tour, error) {

	println("Searching tours for authorID:", authorID)
	println(len(authorID))
	for i, c := range authorID {
		println(i, c, string(c))
	}

	var tours []*model.Tour
	if err := r.db.Where("id_author = ?", authorID).Find(&tours).Error; err != nil {
		return nil, err
	}

	println("Found tours:", len(tours))

	return tours, nil
}

func (r *gormTourRepo) Update(ctx context.Context, t *model.Tour) error {
	if err := r.db.WithContext(ctx).
		Model(&model.Tour{}).
		Where("id = ?", t.ID).
		Updates(map[string]interface{}{
			"name":         t.Name,
			"description":  t.Description,
			"price":        t.Price,
			"published_at": t.PublishedAt,
			"dificulty":    t.Dificulty,
			"status":       t.Status,
			"tags":         t.Tags,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (r *gormTourRepo) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Tour{})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrTourNotFound
	}
	return nil
}
