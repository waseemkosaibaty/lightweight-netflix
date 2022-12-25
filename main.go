package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wkosaibaty/lightweight-netflix/config"
	"github.com/wkosaibaty/lightweight-netflix/controllers"
	"github.com/wkosaibaty/lightweight-netflix/repositories"
	"github.com/wkosaibaty/lightweight-netflix/routes"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoClient *mongo.Client

	userCollection  *mongo.Collection
	movieCollection *mongo.Collection

	userRepository  repositories.UserRepository
	movieRepository repositories.MovieRepository

	authController  controllers.AuthController
	userController  controllers.UserController
	movieController controllers.MovieController

	authRoutes  routes.AuthRoutes
	userRoutes  routes.UserRoutes
	movieRoutes routes.MovieRoutes
)

func init() {
	configuration, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

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
	movieCollection = mongoClient.Database("lightweight_netflix").Collection("movies")

	userRepository = repositories.NewUserRepository(userCollection, ctx)
	movieRepository = repositories.NewMovieRepository(movieCollection, ctx)

	authController = controllers.NewAuthController(userRepository)
	userController = controllers.NewUserController(userRepository)
	movieController = controllers.NewMovieController(movieRepository)

	authRoutes = routes.NewAuthRoutes(authController)
	userRoutes = routes.NewUserRoutes(userController)
	movieRoutes = routes.NewMovieRoutes(movieController)

	server = gin.Default()
}

func main() {
	configuration, _ := config.LoadConfig(".")

	defer mongoClient.Disconnect(ctx)

	server.Static("/public", "./public")

	router := server.Group("/api")
	authRoutes.AddAuthRoutes(router)
	userRoutes.AddUserRoutes(router)
	movieRoutes.AddMovieRoutes(router)

	log.Fatal(server.Run(":" + configuration.Port))
}
