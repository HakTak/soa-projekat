package models

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInternal     = errors.New("internal server error")
)

type FollowRequest struct {
	TargetID string `json:"target_id" binding:"required"`
}

type UserStats struct {
	UserId         string `json:"user_id"`
	FollowersCount int64  `json:"followers_count"`
	FollowingCount int64  `json:"following_count"`
}

type UserRecommendation struct {
	UserID            string `json:"user_id"`
	MutualConnections int64  `json:"mutual_connections"`
}
