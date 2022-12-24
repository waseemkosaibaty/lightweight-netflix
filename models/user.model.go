package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Birthdate time.Time          `json:"birthdate" bson:"birthdate"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
}

type RegisterRequest struct {
	Name      string    `json:"name" bson:"name" binding:"required,min=3,max=255"`
	Birthdate time.Time `json:"birthdate,omitempty" bson:"birthdate,omitempty"`
	Email     string    `json:"email" bson:"email" binding:"required,min=3,max=255"`
	Password  string    `json:"password" bson:"password" binding:"required,min=6,max=255"`
}

type LoginRequest struct {
	Email    string `json:"email" bson:"email" binding:"required,min=3,max=255"`
	Password string `json:"password" bson:"password" binding:"required,min=6,max=255"`
}
