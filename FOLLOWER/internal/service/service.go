package service

import (
	"FOLLOWER/internal/models"
	"FOLLOWER/internal/repository"
	"context"
)

type FollowService interface {
	FollowUser(ctx context.Context, followerID, targetID string) error
	UnfollowUser(ctx context.Context, followerID, targetID string) error
	GetUserStats(ctx context.Context, userID string) (*models.UserStats, error)
	GetFollowers(ctx context.Context, userID string) ([]string, error)
	GetFollowing(ctx context.Context, userID string) ([]string, error)
	GetRecommendations(ctx context.Context, userID string) ([]models.UserRecommendation, error)
	IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error)
}

type followService struct {
	repo repository.SocialRepository
}

func NewFollowService(repo repository.SocialRepository) FollowService {
	return &followService{repo: repo}
}

// FollowUSer implements FollowService.
func (f *followService) FollowUser(ctx context.Context, followerID string, targetID string) error {
	if followerID == targetID {
		return models.ErrInternal
	}
	return f.repo.FollowUser(ctx, followerID, targetID)
}

// GetFollowers implements FollowService.
func (f *followService) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	return f.repo.GetFollowers(ctx, userID)
}

// GetFollowing implements FollowService.
func (f *followService) GetFollowing(ctx context.Context, userID string) ([]string, error) {
	return f.repo.GetFollowing(ctx, userID)
}

// GetRecommendations implements FollowService.
func (f *followService) GetRecommendations(ctx context.Context, userID string) ([]models.UserRecommendation, error) {
	return f.repo.GetRecommendations(ctx, userID)
}

// GetUserStats implements FollowService.
func (f *followService) GetUserStats(ctx context.Context, userID string) (*models.UserStats, error) {
	return f.repo.GetUserStats(ctx, userID)
}

// UnfollowUser implements FollowService.
func (f *followService) UnfollowUser(ctx context.Context, followerID string, targetID string) error {
	return f.repo.UnfollowUser(ctx, followerID, targetID)
}

func (s *followService) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	return s.repo.IsFollowing(ctx, followerID, followeeID)
}
