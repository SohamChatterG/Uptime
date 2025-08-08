package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Check struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UrlID          primitive.ObjectID `bson:"url_id" json:"url_id"`
	UserID         primitive.ObjectID `bson:"user_id" json:"-"`
	CheckedAt      time.Time          `bson:"checked_at" json:"checked_at"`
	WasSuccessful  bool               `bson:"was_successful" json:"was_successful"`
	ResponseTimeMS int64              `bson:"response_time_ms" json:"response_time_ms"`
	StatusCode     int                `bson:"status_code" json:"status_code"`
}
