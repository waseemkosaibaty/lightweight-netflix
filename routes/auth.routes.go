package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wkosaibaty/lightweight-netflix/controllers"
)

type AuthRoutes struct {
	authController controllers.AuthController
}

func NewAuthRoutes(authController controllers.AuthController) AuthRoutes {
	return AuthRoutes{authController}
}

func (authRoutes *AuthRoutes) AddAuthRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/auth")
	router.POST("/", authRoutes.authController.Login)
}
