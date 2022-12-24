package services

import "github.com/wkosaibaty/lightweight-netflix/models"

type UserService interface {
	FindUserByEmail(string) (*models.User, error)
	CreateUser(*models.RegisterRequest) (*models.User, error)
}
