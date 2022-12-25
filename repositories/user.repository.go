package repositories

import "github.com/wkosaibaty/lightweight-netflix/models"

type UserRepository interface {
	FindUserByEmail(string) (*models.User, error)
	CreateUser(*models.RegisterRequest) (*models.User, error)
}
