// service/url_service.go
package service

import (
	"context"
	"errors"
	"time"

	"github.com/SohamChatterG/uptime/model"
	"github.com/SohamChatterG/uptime/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLService struct {
	urlRepo   *repository.URLRepository
	checkRepo *repository.CheckRepository
}

func NewURLService(urlRepo *repository.URLRepository, checkRepo *repository.CheckRepository) *URLService {
	return &URLService{urlRepo: urlRepo, checkRepo: checkRepo}
}

func (s *URLService) CreateURL(ctx context.Context, name, url string, userID primitive.ObjectID) (*model.Url, error) {
	if name == "" || url == "" {
		return nil, errors.New("name and url are required")
	}

	newURL := &model.Url{
		UserID:        userID,
		Name:          name,
		URL:           url,
		Status:        true, // Assume up on creation
		LastCheckedAt: time.Now(),
	}

	return s.urlRepo.Create(ctx, newURL)
}

func (s *URLService) GetURLsForUser(ctx context.Context, userID primitive.ObjectID) ([]model.Url, error) {
	return s.urlRepo.FindByUser(ctx, userID)
}

func (s *URLService) DeleteURL(ctx context.Context, urlIDStr string, userID primitive.ObjectID) error {
	urlID, err := primitive.ObjectIDFromHex(urlIDStr)
	if err != nil {
		return errors.New("invalid url id format")
	}
	return s.urlRepo.DeleteByUser(ctx, urlID, userID)
}

func (s *URLService) GetURLHistory(ctx context.Context, urlIDStr string, userID primitive.ObjectID) ([]model.Check, error) {
	urlID, err := primitive.ObjectIDFromHex(urlIDStr)
	if err != nil {
		return nil, errors.New("invalid url id format")
	}

	// First, verify the user owns this URL before returning its history
	_, err = s.urlRepo.FindByIDAndUser(ctx, urlID, userID)
	if err != nil {
		return nil, errors.New("url not found or access denied")
	}

	// Now fetch the history
	return s.checkRepo.GetHistoryForURL(ctx, urlIDStr, 100) // Limit to last 100 checks
}
