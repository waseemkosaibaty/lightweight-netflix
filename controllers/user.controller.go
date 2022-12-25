package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wkosaibaty/lightweight-netflix/models"
	"github.com/wkosaibaty/lightweight-netflix/repositories"
	"github.com/wkosaibaty/lightweight-netflix/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	userRepository repositories.UserRepository
}

func NewUserController(userRepository repositories.UserRepository) UserController {
	return UserController{userRepository}
}

func (controller *UserController) CreateUser(ctx *gin.Context) {
	var request *models.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := controller.userRepository.FindUserByEmail(request.Email)
	if user != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email already exists"})
		return
	}

	if err != mongo.ErrNoDocuments {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	user, err = controller.userRepository.CreateUser(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	access_token, err := utils.CreateJWT(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"access_token": access_token})
}
