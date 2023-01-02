package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Recipe struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name"`
	Tags         []string           `json:"tags"`
	Ingredients  []string           `json:"ingredients"`
	Instructions []string           `json:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt"`
}
