package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WatchedMovie struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	MovieId primitive.ObjectID `json:"movieId" bson:"movieId"`
	Movie   Movie              `json:"movie,omitempty" bson:"movie,omitempty"`
	UserId  primitive.ObjectID `json:"userId" bson:"userId"`
	Rate    int                `json:"rate" bson:"rate"`
	Review  string             `json:"review,omitempty" bson:"review,omitempty"`
}

type CreateWatchedMovieRequest struct {
	MovieId primitive.ObjectID `json:"movieId" bson:"movieId" binding:"required"`
	UserId  primitive.ObjectID `json:"userId" bson:"userId"`
	Rate    int                `json:"rate" bson:"rate"`
}

type RateWatchedMovieRequest struct {
	Rate   int    `json:"rate" bson:"rate" binding:"required,min=1,max=5"`
	Review string `json:"review,omitempty" bson:"review,omitempty" binding:"max=1024"`
}
