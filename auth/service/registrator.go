package service

import (
	"auth/hasher"
	"auth/models"
	"auth/repo"
	"context"
)

type Registrator interface {
	Register(ctx context.Context, username string, password string, scopes []string) error
}

type registrator struct {
	repo   repo.Repo
	hasher hasher.Hash
}

func NewRegistrator(repo repo.Repo, hasher hasher.Hash) Registrator {
	return &registrator{
		repo:   repo,
		hasher: hasher,
	}
}

func (registratorService *registrator) Register(
	ctx context.Context,
	username string,
	password string,
	scopes []string,
) error {
	hashedPassword, err := registratorService.hasher.HashPassword(password)
	if err != nil {
		return err
	}

	user := models.User{
		Username:       username,
		HashedPassword: hashedPassword,
		Scopes:         scopes,
	}

	err = registratorService.repo.InsertUser(ctx, user)
	return err
}
