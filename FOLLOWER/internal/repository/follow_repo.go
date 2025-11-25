package repository

import (
	"FOLLOWER/internal/models"
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type SocialRepository interface {
	EnsureConstraints(ctx context.Context) error
	FollowUser(ctx context.Context, followerID, targetID string) error
	UnfollowUser(ctx context.Context, followerID, targetID string) error
	GetUserStats(ctx context.Context, userID string) (*models.UserStats, error)
	GetFollowers(ctx context.Context, userID string) ([]string, error)
	GetFollowing(ctx context.Context, userID string) ([]string, error)
	GetRecommendations(ctx context.Context, userID string) ([]models.UserRecommendation, error)
	IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error)
}

type followRepository struct {
	Driver neo4j.DriverWithContext
}

func NewfollowRepository(driver neo4j.DriverWithContext) *followRepository {
	return &followRepository{Driver: driver}
}

func (r *followRepository) EnsureConstraints(ctx context.Context) error {
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := `
        CREATE CONSTRAINT user_id_unique IF NOT EXISTS 
        FOR (u:User) 
        REQUIRE u.userId IS UNIQUE
    `

	_, err := session.Run(ctx, query, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *followRepository) FollowUser(ctx context.Context, followerID, followeeID string) error {
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := `
			MERGE (u1:User {userId: $followerID})
			MERGE (u2:User {userId: $followeeID})
			MERGE (u1)-[:FOLLOWS]->(u2)
	`
	_, err := session.Run(ctx, query, map[string]any{"followerID": followerID, "followeeID": followeeID})

	return err
}

func (r *followRepository) UnfollowUser(ctx context.Context, followerID, followeeID string) error {
	sesion := r.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer sesion.Close(ctx)

	query := `
			MATCH (u1:User {userId: $followerID})-[r:FOLLOWS]->(u2:User {userId: $followeeID}) DELETE r
	`
	_, err := sesion.Run(ctx, query, map[string]any{"followerID": followerID, "followeeID": followeeID})

	return err
}

func (r *followRepository) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	res, err := session.Run(ctx, "MATCH (me:User {userId: $id})<-[:FOLLOWS]-(other) RETURN other.userId as id", map[string]any{"id": userID})
	if err != nil {
		return nil, err
	}
	var ids []string
	for res.Next(ctx) {
		val, _ := res.Record().Get("id")
		ids = append(ids, val.(string))
	}
	return ids, nil
}

func (r *followRepository) GetFollowing(ctx context.Context, userID string) ([]string, error) {
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	res, err := session.Run(ctx, "MATCH (me:User {userId: $id})-[:FOLLOWS]->(other) RETURN other.userId as id", map[string]any{"id": userID})
	if err != nil {
		return nil, err
	}
	var ids []string
	for res.Next(ctx) {
		val, _ := res.Record().Get("id")
		ids = append(ids, val.(string))
	}
	return ids, nil
}

func (r *followRepository) GetUserStats(ctx context.Context, userID string) (*models.UserStats, error) {
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := `
		MATCH (u:User {userId: $id})
		RETURN count((u)<-[:FOLLOWS]-()) AS followers, count((u)-[:FOLLOWS]->()) AS following
	`
	result, err := session.Run(ctx, query, map[string]any{"id": userID})
	if err != nil {
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		f1, _ := record.Get("followers")
		f2, _ := record.Get("following")
		return &models.UserStats{UserId: userID, FollowersCount: f1.(int64), FollowingCount: f2.(int64)}, nil
	}

	return &models.UserStats{UserId: userID, FollowersCount: 0, FollowingCount: 0}, nil
}

func (r *followRepository) GetRecommendations(ctx context.Context, userID string) ([]models.UserRecommendation, error) {
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := `
		MATCH (me:User {userId: $id})-[:FOLLOWS]->(friend)-[:FOLLOWS]->(fof:User)
		WHERE NOT (me)-[:FOLLOWS]->(fof) AND me <> fof
		RETURN fof.userId AS id, count(friend) AS mutual
		ORDER BY mutual DESC LIMIT 10
	`
	result, err := session.Run(ctx, query, map[string]any{"id": userID})
	if err != nil {
		return nil, err
	}

	var recs []models.UserRecommendation
	for result.Next(ctx) {
		rec := result.Record()
		id, _ := rec.Get("id")
		mutual, _ := rec.Get("mutual")
		recs = append(recs, models.UserRecommendation{UserID: id.(string), MutualConnections: mutual.(int64)})
	}
	return recs, nil
}

func (r *followRepository) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	// Cypher: Look for the pattern. Return true if found, false if not.
	query := `
        MATCH (a:User {userId: $follower})-[r:FOLLOWS]->(b:User {userId: $followee})
        RETURN count(r) > 0 AS exists
    `

	params := map[string]any{
		"follower": followerID,
		"followee": followeeID,
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return false, err
	}

	if result.Next(ctx) {
		val, _ := result.Record().Get("exists")
		return val.(bool), nil
	}

	return false, nil
}
