package repo

import (
	"auth/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DatabaseName    string = "mx-auth"
	UsersCollection string = "users"
)

type mongoRepo struct {
	client *mongo.Client
}

func NewMongoRepo(client *mongo.Client) Repo {
	return &mongoRepo{
		client: client,
	}
}

func (repo *mongoRepo) GetUser(ctx context.Context, username string) (models.User, error) {
	filter := bson.M{
		"username": username,
	}
	user := models.User{}
	result := repo.client.Database(DatabaseName).Collection(UsersCollection).FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		return user, err
	}

	err := result.Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (repo *mongoRepo) InsertUser(ctx context.Context, user models.User) error {
	_, err := repo.client.Database(DatabaseName).Collection(UsersCollection).InsertOne(ctx, user)
	return err
}
