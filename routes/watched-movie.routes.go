package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wkosaibaty/lightweight-netflix/controllers"
	"github.com/wkosaibaty/lightweight-netflix/middleware"
)

type WatchedMovieRoutes struct {
	controller controllers.WatchedMovieController
}

func NewWatchedMovieRoutes(movieController controllers.WatchedMovieController) WatchedMovieRoutes {
	return WatchedMovieRoutes{movieController}
}

func (routes *WatchedMovieRoutes) AddWatchedMovieRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/watched-movies")
	router.GET("/", middleware.AuthMiddleware(), routes.controller.FindUserWatchedMovies)
	router.POST("/", middleware.AuthMiddleware(), routes.controller.CreateWatchedMovie)
	router.PUT("/:movieId", middleware.AuthMiddleware(), routes.controller.RateWatchedMovie)
}
