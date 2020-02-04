package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Assignment is the object model for the Assignment collection
type Assignment struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Type        string             `json:"type" bson:"type"`
	Duration    int                `json:"duration" bson:"duration"`
	Status      string             `json:"status" bson:"status"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	Tags        []string           `json:"tags" bson:"tags"`
}
