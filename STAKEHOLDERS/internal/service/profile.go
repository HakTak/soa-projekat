package service

import (
	"context"
	"errors"
	"stakeholders/internal/model"
	"stakeholders/internal/repository"
)

var (
	ErrProfileExists      = errors.New("profile already exists")
	ErrCannotBlockAdmin   = errors.New("cannot block an admin user")
	ErrMissingProfileData = errors.New("missing requierd profile data")
)

type ProfileService interface {
	CreateProfile(ctx context.Context, p *model.Profile) (*model.Profile, error)
	GetProfile(ctx context.Context, userID string) (*model.Profile, error)
	UpdateProfile(ctx context.Context, p *model.Profile) (*model.Profile, error)
	BlockUser(ctx context.Context, userID string, adminUserID string) error
}

type profileService struct {
	repo repository.ProfileRepository
}

func NewProfileService(repo repository.ProfileRepository) ProfileService {
	return &profileService{repo: repo}
}

func (s *profileService) BlockUser(ctx context.Context, userID string, adminUserID string) error {
	p, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if p.Role == model.RoleAdmin {
		return ErrCannotBlockAdmin
	}
	return s.repo.Block(ctx, userID)
}

// CreateProfile implements ProfileService.
func (s *profileService) CreateProfile(ctx context.Context, p *model.Profile) (*model.Profile, error) {
	if p.UserID == "" || p.Role == "" {
		return nil, ErrMissingProfileData
	}

	if err := s.repo.Create(ctx, p); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return nil, ErrProfileExists
		}
		return nil, err
	}
	return p, nil
}

// GetProfile implements ProfileService.
func (s *profileService) GetProfile(ctx context.Context, userID string) (*model.Profile, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// UpdateProfile implements ProfileService.
func (s *profileService) UpdateProfile(ctx context.Context, p *model.Profile) (*model.Profile, error) {
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return s.repo.GetByUserID(ctx, p.UserID)
}
