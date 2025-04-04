package repo

import (
	"auth/models"
	"context"
)

type Repo interface {
	GetHashedPasswordByUsername(ctx context.Context, username string) (string, error)
	InsertUser(ctx context.Context, user models.User) error
}
