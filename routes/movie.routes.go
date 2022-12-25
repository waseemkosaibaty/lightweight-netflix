package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wkosaibaty/lightweight-netflix/controllers"
	"github.com/wkosaibaty/lightweight-netflix/middleware"
)

type MovieRoutes struct {
	movieController controllers.MovieController
}

func NewMovieRoutes(movieController controllers.MovieController) MovieRoutes {
	return MovieRoutes{movieController}
}

func (movieRoutes *MovieRoutes) AddMovieRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/movies")
	router.GET("/", movieRoutes.movieController.FindAllMovies)
	router.GET("/:id", movieRoutes.movieController.FindMovieById)
	router.POST("/", middleware.AuthMiddleware(), movieRoutes.movieController.CreateMovie)
	router.PUT("/:id", middleware.AuthMiddleware(), movieRoutes.movieController.UpdateMovie)
	router.DELETE("/:id", middleware.AuthMiddleware(), movieRoutes.movieController.DeleteMovie)
}
