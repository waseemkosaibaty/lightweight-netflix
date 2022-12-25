package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wkosaibaty/lightweight-netflix/models"
	"github.com/wkosaibaty/lightweight-netflix/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WatchedMovieController struct {
	watchedMovieRepository repositories.WatchedMovieRepository
	movieRepository        repositories.MovieRepository
}

func NewWatchedMovieController(watchedMovieRepository repositories.WatchedMovieRepository, movieRepository repositories.MovieRepository) WatchedMovieController {
	return WatchedMovieController{watchedMovieRepository, movieRepository}
}

func (controller *WatchedMovieController) FindUserWatchedMovies(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(string)

	watchedMovies, err := controller.watchedMovieRepository.FindAllWatchedMovies(userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, watchedMovies)
}

func (controller *WatchedMovieController) CreateWatchedMovie(ctx *gin.Context) {
	var request *models.CreateWatchedMovieRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userId := ctx.MustGet("userId").(string)
	request.UserId, _ = primitive.ObjectIDFromHex(userId)

	_, err := controller.movieRepository.FindMovieById(request.MovieId.Hex())
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		}
		return
	}

	watchedMovie, err := controller.watchedMovieRepository.FindWatchedMovie(userId, request.MovieId.Hex())
	if watchedMovie != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "This movie was added to your watch list before"})
		return
	}

	if err != mongo.ErrNoDocuments {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	watchedMovie, err = controller.watchedMovieRepository.CreateWatchedMovie(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, watchedMovie)
}

func (controller *WatchedMovieController) RateWatchedMovie(ctx *gin.Context) {
	movieId := ctx.Param("movieId")
	userId := ctx.MustGet("userId").(string)

	var request *models.RateWatchedMovieRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	watchedMovie, err := controller.watchedMovieRepository.RateWatchedMovie(userId, movieId, request)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, watchedMovie)
}
