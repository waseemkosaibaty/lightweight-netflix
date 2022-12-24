package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wkosaibaty/lightweight-netflix/models"
	"github.com/wkosaibaty/lightweight-netflix/services"
	"github.com/wkosaibaty/lightweight-netflix/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	userService services.UserService
	ctx         context.Context
	collection  *mongo.Collection
}

func NewUserController(userService services.UserService, ctx context.Context, collection *mongo.Collection) UserController {
	return UserController{userService, ctx, collection}
}

func (controller *UserController) CreateUser(ctx *gin.Context) {
	var request *models.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := controller.userService.CreateUser(request)
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
