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
	// Additional methods can be added later.
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
