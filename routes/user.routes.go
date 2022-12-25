package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wkosaibaty/lightweight-netflix/controllers"
)

type UserRoutes struct {
	userController controllers.UserController
}

func NewUserRoutes(userController controllers.UserController) UserRoutes {
	return UserRoutes{userController}
}

func (userRoutes *UserRoutes) AddUserRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/users")
	router.POST("/", userRoutes.userController.CreateUser)
}
