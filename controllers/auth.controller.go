package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wkosaibaty/lightweight-netflix/models"
	"github.com/wkosaibaty/lightweight-netflix/services"
	"github.com/wkosaibaty/lightweight-netflix/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	userService services.UserService
	ctx         context.Context
	collection  *mongo.Collection
}

func NewAuthController(userService services.UserService, ctx context.Context, collection *mongo.Collection) AuthController {
	return AuthController{userService, ctx, collection}
}

func (controller *AuthController) Login(ctx *gin.Context) {
	errorMessage := gin.H{"message": "Invalid email or password"}

	var request *models.LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := controller.userService.FindUserByEmail(request.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusBadRequest, errorMessage)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage)
		return
	}

	access_token, err := utils.CreateJWT(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": access_token})
}
