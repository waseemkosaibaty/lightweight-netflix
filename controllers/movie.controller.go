package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wkosaibaty/lightweight-netflix/models"
	"github.com/wkosaibaty/lightweight-netflix/repositories"
	"github.com/wkosaibaty/lightweight-netflix/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieController struct {
	movieRepository repositories.MovieRepository
}

func NewMovieController(movieRepository repositories.MovieRepository) MovieController {
	return MovieController{movieRepository}
}

func (controller *MovieController) FindAllMovies(ctx *gin.Context) {
	sortBy := strings.ToLower(ctx.DefaultQuery("sortBy", ""))
	sortTypeString := strings.ToLower(ctx.DefaultQuery("sortType", "asc"))

	sortType := 1
	if sortTypeString == "desc" {
		sortType = -1
	}

	movies, err := controller.movieRepository.FindAllMovies(sortBy, sortType)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, movies)
}

func (controller *MovieController) FindMovieById(ctx *gin.Context) {
	id := ctx.Param("id")

	movie, err := controller.movieRepository.FindMovieById(id)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, movie)
}

func (controller *MovieController) CreateMovie(ctx *gin.Context) {
	var request *models.CreateMovieRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userId := ctx.MustGet("userId").(string)
	request.UserId, _ = primitive.ObjectIDFromHex(userId)

	if request.Cover != "" {
		coverUrl, err := utils.UploadImage(request.Cover)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		request.Cover = coverUrl
	}

	movie, err := controller.movieRepository.CreateMovie(request)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, movie)
}

func (controller *MovieController) UpdateMovie(ctx *gin.Context) {
	id := ctx.Param("id")
	userId := ctx.MustGet("userId").(string)
	userObjectId, _ := primitive.ObjectIDFromHex(userId)

	movie, err := controller.movieRepository.FindMovieById(id)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		}
		return
	}

	if movie.UserId != userObjectId {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "You do not have access to modify this movie"})
		return
	}

	var request *models.UpdateMovieRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if request.Cover != "" {
		coverUrl, err := utils.UploadImage(request.Cover)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			return
		}
		request.Cover = coverUrl
	}

	movie, err = controller.movieRepository.UpdateMovie(id, request)
	if err != nil {
		if err != nil {
			if strings.Contains(err.Error(), "not exist") {
				ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			} else {
				ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			}
			return
		}
	}

	ctx.JSON(http.StatusOK, movie)
}

func (controller *MovieController) DeleteMovie(ctx *gin.Context) {
	id := ctx.Param("id")
	userId := ctx.MustGet("userId").(string)
	userObjectId, _ := primitive.ObjectIDFromHex(userId)

	movie, err := controller.movieRepository.FindMovieById(id)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		}
		return
	}

	if movie.UserId != userObjectId {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "You do not have access to delete this movie"})
		return
	}

	err = controller.movieRepository.DeleteMovie(id)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, id)
}
