package repositories

import "github.com/wkosaibaty/lightweight-netflix/models"

type MovieRepository interface {
	FindAllMovies(string, int) ([]*models.Movie, error)
	FindMovieById(string) (*models.Movie, error)
	CreateMovie(*models.CreateMovieRequest) (*models.Movie, error)
	UpdateMovie(string, *models.UpdateMovieRequest) (*models.Movie, error)
	DeleteMovie(string) error
}
