package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"stakeholders/internal/model"
)

var (
	ErrNotFound      = errors.New("profile not found")
	ErrAlreadyExists = errors.New("profile already exists")
)

type ProfileRepository interface {
	Create(ctx context.Context, p *model.Profile) error
	GetByUserID(ctx context.Context, userID string) (*model.Profile, error)
	Update(ctx context.Context, p *model.Profile) error
	Block(ctx context.Context, userID string) error
}

type gormProfileRepo struct {
	db *gorm.DB
}

func NewGormProfileRepo(db *gorm.DB) ProfileRepository {
	return &gormProfileRepo{db: db}
}

func (r *gormProfileRepo) Create(ctx context.Context, p *model.Profile) error {
	if err := r.db.WithContext(ctx).Create(p).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *gormProfileRepo) GetByUserID(ctx context.Context, userID string) (*model.Profile, error) {
	var p model.Profile
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *gormProfileRepo) Update(ctx context.Context, p *model.Profile) error {
	if err := r.db.WithContext(ctx).Model(&model.Profile{}).Where("user_id = ?", p.UserID).Updates(map[string]interface{}{
		"first_name":      p.FirstName,
		"last_name":       p.LastName,
		"profile_picture": p.ProfilePicture,
		"biography":       p.Biography,
		"motto":           p.Motto,
		"role":            p.Role,
		"is_blocked":      p.IsBlocked,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *gormProfileRepo) Block(ctx context.Context, userID string) error {
	res := r.db.WithContext(ctx).Model(&model.Profile{}).Where("user_id = ?", userID).Update("is_blocked", true)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
