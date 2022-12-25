package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Date        time.Time          `json:"date" bson:"date"`
	Cover       string             `json:"cover,omitempty" bson:"cover,omitempty"`
	UserId      primitive.ObjectID `json:"user" bson:"user"`
}

type CreateMovieRequest struct {
	Name        string             `json:"name" bson:"name" binding:"required,min=3,max=255"`
	Description string             `json:"description,omitempty" bson:"description,omitempty" binding:"max=1024"`
	Date        time.Time          `json:"date" bson:"date" binding:"required"`
	Cover       string             `json:"cover,omitempty" bson:"cover,omitempty"`
	UserId      primitive.ObjectID `json:"user,omitempty" bson:"user,omitempty"`
}

type UpdateMovieRequest struct {
	Name        string    `json:"name" bson:"name" binding:"required,min=3,max=255"`
	Description string    `json:"description,omitempty" bson:"description,omitempty" binding:"max=1024"`
	Date        time.Time `json:"date" bson:"date" binding:"required"`
	Cover       string    `json:"cover,omitempty" bson:"cover,omitempty"`
}
