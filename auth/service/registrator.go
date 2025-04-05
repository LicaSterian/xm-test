package service

import (
	"auth/hasher"
	"auth/models"
	"auth/repo"
	"context"
	"errors"
)

var ErrHashingPassword = errors.New("error hashing password")
var ErrInsertingUser = errors.New("error inserting user in the database")

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
		return errors.Join(ErrHashingPassword, err)
	}

	user := models.User{
		Username:       username,
		HashedPassword: hashedPassword,
		Scopes:         scopes,
	}

	err = registratorService.repo.InsertUser(ctx, user)
	if err != nil {
		return errors.Join(ErrInsertingUser, err)
	}
	return nil
}
