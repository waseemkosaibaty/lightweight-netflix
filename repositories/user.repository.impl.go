package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/wkosaibaty/lightweight-netflix/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserRepository(collection *mongo.Collection, ctx context.Context) UserRepository {
	return &UserRepositoryImpl{collection, ctx}
}

func (repository *UserRepositoryImpl) FindUserByEmail(email string) (*models.User, error) {
	var user *models.User

	query := bson.M{"email": strings.ToLower(strings.Trim(email, " "))}
	if err := repository.collection.FindOne(repository.ctx, query).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (repository *UserRepositoryImpl) CreateUser(request *models.RegisterRequest) (*models.User, error) {
	request.Email = strings.ToLower(strings.Trim(request.Email, " "))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("Could not hash password")
	}
	request.Password = string(hashedPassword)

	result, err := repository.collection.InsertOne(repository.ctx, &request)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("Email already exists")
		}
		return nil, err
	}

	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true)}
	if _, err := repository.collection.Indexes().CreateOne(repository.ctx, index); err != nil {
		return nil, errors.New("Could not create index for email")
	}

	var user *models.User
	err = repository.collection.FindOne(repository.ctx, bson.M{"_id": result.InsertedID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
