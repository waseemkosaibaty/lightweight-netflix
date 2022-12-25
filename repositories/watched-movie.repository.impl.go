package repositories

import (
	"context"
	"errors"

	"github.com/wkosaibaty/lightweight-netflix/models"
	"github.com/wkosaibaty/lightweight-netflix/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WatchedMovieRepositoryImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewWatchedMovieRepository(collection *mongo.Collection, ctx context.Context) WatchedMovieRepository {
	return &WatchedMovieRepositoryImpl{collection, ctx}
}

func (repository *WatchedMovieRepositoryImpl) FindAllWatchedMovies(userId string) ([]*models.WatchedMovie, error) {
	userObjectId, _ := primitive.ObjectIDFromHex(userId)

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "userId", Value: userObjectId}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "movies"}, {Key: "localField", Value: "movieId"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "movie"}}}}
	unwindStage := bson.D{{Key: "$unwind", Value: "$movie"}}

	cursor, err := repository.collection.Aggregate(repository.ctx, mongo.Pipeline{matchStage, lookupStage, unwindStage})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(repository.ctx)

	var watchedMovies []*models.WatchedMovie

	for cursor.Next(repository.ctx) {
		watchedMovie := &models.WatchedMovie{}
		err := cursor.Decode(watchedMovie)
		if err != nil {
			return nil, err
		}
		watchedMovies = append(watchedMovies, watchedMovie)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(watchedMovies) == 0 {
		return []*models.WatchedMovie{}, nil
	}
	return watchedMovies, nil
}

func (repository *WatchedMovieRepositoryImpl) FindWatchedMovie(userId string, movieId string) (*models.WatchedMovie, error) {
	userObjectId, _ := primitive.ObjectIDFromHex(userId)
	movieObjectId, _ := primitive.ObjectIDFromHex(movieId)
	var watchedMovie *models.WatchedMovie

	err := repository.collection.FindOne(repository.ctx, bson.M{"userId": userObjectId, "movieId": movieObjectId}).Decode(&watchedMovie)
	if err != nil {
		return nil, err
	}

	return watchedMovie, nil
}

func (repository *WatchedMovieRepositoryImpl) CreateWatchedMovie(request *models.CreateWatchedMovieRequest) (*models.WatchedMovie, error) {
	request.Rate = 0

	res, err := repository.collection.InsertOne(repository.ctx, request)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("You already watched this movie")
		}
		return nil, err
	}

	var watchedMovie *models.WatchedMovie
	if err = repository.collection.FindOne(repository.ctx, bson.M{"_id": res.InsertedID}).Decode(&watchedMovie); err != nil {
		return nil, err
	}

	return watchedMovie, nil
}

func (repository *WatchedMovieRepositoryImpl) RateWatchedMovie(userId string, movieId string, request *models.RateWatchedMovieRequest) (*models.WatchedMovie, error) {
	document, err := utils.ToDocument(request)
	if err != nil {
		return nil, err
	}

	userObjectId, _ := primitive.ObjectIDFromHex(userId)
	movieObjectId, _ := primitive.ObjectIDFromHex(movieId)
	query := bson.D{{Key: "userId", Value: userObjectId}, {Key: "movieId", Value: movieObjectId}, {Key: "rate", Value: 0}}
	update := bson.D{{Key: "$set", Value: document}}
	result := repository.collection.FindOneAndUpdate(repository.ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var watchedMovie *models.WatchedMovie
	if err := result.Decode(&watchedMovie); err != nil {
		return nil, errors.New("Cannot rate the specified movie. Either it doesn't exist or you rated this movie before")
	}

	return watchedMovie, nil
}
