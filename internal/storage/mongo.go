package storage

import (
	"context"
	"errors"

	"github.com/checkly-go/checkly/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository is defined for future user management.
type UserRepository interface {
	// Future methods for user management.
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

// CheckRepository provides DB operations for website checks.
type CheckRepository interface {
	CreateCheck(ctx context.Context, check *models.WebsiteCheck) error
	GetCheck(ctx context.Context, id primitive.ObjectID) (*models.WebsiteCheck, error)
	GetLeaderboard(ctx context.Context, limit int) ([]models.LeaderboardEntry, error)
	GetAllChecks(ctx context.Context) ([]models.WebsiteCheck, error)
}

type checkRepository struct {
	collection *mongo.Collection
}

func NewCheckRepository(db *mongo.Database) CheckRepository {
	return &checkRepository{
		collection: db.Collection("website_checks"),
	}
}

func (r *checkRepository) CreateCheck(ctx context.Context, check *models.WebsiteCheck) error {
	res, err := r.collection.InsertOne(ctx, check)
	if err != nil {
		return err
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("failed to assert inserted id as ObjectID")
	}
	check.ID = id
	return nil
}

func (r *checkRepository) GetCheck(ctx context.Context, id primitive.ObjectID) (*models.WebsiteCheck, error) {
	var result models.WebsiteCheck
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *checkRepository) GetLeaderboard(ctx context.Context, limit int) ([]models.LeaderboardEntry, error) {
	if limit <= 0 {
		limit = 10
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"report":              bson.M{"$ne": nil},
				"report.overallscore": bson.M{"$ne": nil, "$exists": true},
			},
		},
		{
			"$sort": bson.M{
				"report.overallscore": -1,
				"created_at":          -1,
			},
		},
		{
			"$group": bson.M{
				"_id":          "$url",
				"url":          bson.M{"$first": "$url"},
				"overallscore": bson.M{"$first": "$report.overallscore"},
				"timestamp":    bson.M{"$first": "$created_at"},
			},
		},
		{
			"$sort": bson.M{
				"overallscore": -1,
				"timestamp":    -1,
			},
		},
		{
			"$limit": limit,
		},
		{
			"$project": bson.M{
				"url":          1,
				"overallscore": 1,
				"timestamp":    1,
				"_id":          0,
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.LeaderboardEntry
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *checkRepository) GetAllChecks(ctx context.Context) ([]models.WebsiteCheck, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.WebsiteCheck
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
