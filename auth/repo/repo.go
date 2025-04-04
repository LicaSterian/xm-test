package repo

import (
	"auth/models"
	"context"
)

type Repo interface {
	GetUser(ctx context.Context, username string) (models.User, error)
	InsertUser(ctx context.Context, user models.User) error
}
