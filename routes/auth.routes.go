package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wkosaibaty/lightweight-netflix/controllers"
	"github.com/wkosaibaty/lightweight-netflix/services"
)

type AuthRoutes struct {
	authController controllers.AuthController
}

func NewAuthRoutes(authController controllers.AuthController) AuthRoutes {
	return AuthRoutes{authController}
}

func (authRoutes *AuthRoutes) AddAuthRoutes(rg *gin.RouterGroup, userService services.UserService) {
	router := rg.Group("/auth")
	router.POST("/", authRoutes.authController.Login)
}
