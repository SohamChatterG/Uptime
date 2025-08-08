package repository

import (
	"context"

	"github.com/SohamChatterG/uptime/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive" // <-- 1. ADD THIS IMPORT
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CheckRepository struct {
	collection *mongo.Collection
}

func NewCheckRepository(db *mongo.Database) *CheckRepository {
	return &CheckRepository{
		collection: db.Collection("checks"),
	}
}

func (r *CheckRepository) Create(ctx context.Context, check *model.Check) error {
	_, err := r.collection.InsertOne(ctx, check)
	return err
}

func (r *CheckRepository) GetHistoryForURL(ctx context.Context, urlID string, limit int64) ([]model.Check, error) {
	// 2. CHANGE THE FUNCTION CALL HERE
	objID, err := primitive.ObjectIDFromHex(urlID)
	if err != nil {
		return nil, err
	}

	opts := options.Find().SetSort(bson.D{{"checked_at", -1}}).SetLimit(limit)
	cursor, err := r.collection.Find(ctx, bson.M{"url_id": objID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var checks []model.Check
	if err = cursor.All(ctx, &checks); err != nil {
		return nil, err
	}
	return checks, nil
}
