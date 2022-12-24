package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wkosaibaty/lightweight-netflix/config"
	"github.com/wkosaibaty/lightweight-netflix/controllers"
	"github.com/wkosaibaty/lightweight-netflix/routes"
	"github.com/wkosaibaty/lightweight-netflix/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoClient *mongo.Client

	userCollection *mongo.Collection

	userService services.UserService

	authController controllers.AuthController
	userController controllers.UserController

	authRoutes routes.AuthRoutes
	userRoutes routes.UserRoutes
)

func init() {
	configuration, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	server = gin.Default()

	ctx = context.TODO()

	mongoconn := options.Client().ApplyURI(configuration.MongodbUri)
	mongoClient, err := mongo.Connect(ctx, mongoconn)
	if err != nil {
		panic(err)
	}
	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB...")

	userCollection = mongoClient.Database("lightweight_netflix").Collection("users")

	userService = services.NewUserService(userCollection, ctx)

	authController = controllers.NewAuthController(userService, ctx, userCollection)
	userController = controllers.NewUserController(userService, ctx, userCollection)

	authRoutes = routes.NewAuthRoutes(authController)
	userRoutes = routes.NewUserRoutes(userController)
}

func main() {
	configuration, _ := config.LoadConfig(".")

	defer mongoClient.Disconnect(ctx)

	router := server.Group("/api")
	authRoutes.AddAuthRoutes(router, userService)
	userRoutes.AddUserRoutes(router, userService)

	log.Fatal(server.Run(":" + configuration.Port))
}
