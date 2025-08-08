package repository

import (
	"context"
	"time"

	"github.com/SohamChatterG/uptime/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type URLRepository struct {
	collection *mongo.Collection
}

func NewURLRepository(db *mongo.Database) *URLRepository {
	return &URLRepository{
		collection: db.Collection("urls"),
	}
}

func (r *URLRepository) Create(ctx context.Context, url *model.Url) (*model.Url, error) {
	res, err := r.collection.InsertOne(ctx, url)
	if err != nil {
		return nil, err
	}
	url.ID = res.InsertedID.(primitive.ObjectID)
	return url, nil
}

func (r *URLRepository) FindByUser(ctx context.Context, userID primitive.ObjectID) ([]model.Url, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var urls []model.Url
	if err = cursor.All(ctx, &urls); err != nil {
		return nil, err
	}
	return urls, nil
}

func (r *URLRepository) FindByIDAndUser(ctx context.Context, urlID, userID primitive.ObjectID) (*model.Url, error) {
	var url model.Url
	filter := bson.M{"_id": urlID, "user_id": userID}
	err := r.collection.FindOne(ctx, filter).Decode(&url)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *URLRepository) DeleteByUser(ctx context.Context, urlID, userID primitive.ObjectID) error {
	filter := bson.M{"_id": urlID, "user_id": userID}
	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *URLRepository) GetAllActive(ctx context.Context) ([]model.Url, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}) // Add a filter like {"is_paused": false} in the future
	if err != nil {
		return nil, err
	}
	var urls []model.Url
	if err = cursor.All(ctx, &urls); err != nil {
		return nil, err
	}
	return urls, nil
}

func (r *URLRepository) UpdateStatus(ctx context.Context, urlID primitive.ObjectID, status bool) error {
	filter := bson.M{"_id": urlID}
	update := bson.M{"$set": bson.M{"status": status, "last_checked_at": time.Now()}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
