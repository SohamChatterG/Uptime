// model/url.go
package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Url struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"-"`
	Name          string             `bson:"name" json:"name"`
	URL           string             `bson:"url" json:"url"`
	Status        bool               `bson:"status" json:"status"`
	LastCheckedAt time.Time          `bson:"last_checked_at" json:"last_checked_at"`
}
