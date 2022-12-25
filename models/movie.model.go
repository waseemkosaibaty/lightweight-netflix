package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Date        time.Time          `json:"date" bson:"date"`
	Cover       string             `json:"cover,omitempty" bson:"cover,omitempty"`
	UserId      primitive.ObjectID `json:"userId" bson:"userId"`
	Rate        float64            `json:"rate,omitempty" bson:"rate,omitempty"`
}

type CreateMovieRequest struct {
	Name        string             `json:"name" bson:"name" binding:"required,min=3,max=255"`
	Description string             `json:"description,omitempty" bson:"description,omitempty" binding:"max=1024"`
	Date        time.Time          `json:"date" bson:"date" binding:"required"`
	Cover       string             `json:"cover,omitempty" bson:"cover,omitempty"`
	UserId      primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
}

type UpdateMovieRequest struct {
	Name        string    `json:"name" bson:"name" binding:"required,min=3,max=255"`
	Description string    `json:"description,omitempty" bson:"description,omitempty" binding:"max=10240"`
	Date        time.Time `json:"date" bson:"date" binding:"required"`
	Cover       string    `json:"cover,omitempty" bson:"cover,omitempty"`
}
