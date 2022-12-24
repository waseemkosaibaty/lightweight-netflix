package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wkosaibaty/lightweight-netflix/controllers"
	"github.com/wkosaibaty/lightweight-netflix/services"
)

type UserRoutes struct {
	userController controllers.UserController
}

func NewUserRoutes(userController controllers.UserController) UserRoutes {
	return UserRoutes{userController}
}

func (userRoutes *UserRoutes) AddUserRoutes(rg *gin.RouterGroup, userService services.UserService) {
	router := rg.Group("/users")
	router.POST("/", userRoutes.userController.CreateUser)
}
