package repositories

import "github.com/wkosaibaty/lightweight-netflix/models"

type WatchedMovieRepository interface {
	FindAllWatchedMovies(string) ([]*models.WatchedMovie, error)
	FindWatchedMovie(string, string) (*models.WatchedMovie, error)
	CreateWatchedMovie(*models.CreateWatchedMovieRequest) (*models.WatchedMovie, error)
	RateWatchedMovie(string, string, *models.RateWatchedMovieRequest) (*models.WatchedMovie, error)
}
