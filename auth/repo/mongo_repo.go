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

func (repo *mongoRepo) GetHashedPasswordByUsername(ctx context.Context, username string) (string, error) {
	filter := bson.M{
		"username": username,
	}
	result := repo.client.Database(DatabaseName).Collection(UsersCollection).FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		return "", err
	}
	var user models.User
	err := result.Decode(&user)
	if err != nil {
		return "", err
	}
	return user.HashedPassword, nil
}
