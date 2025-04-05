package service

import (
	"auth/hasher"
	"auth/repo"
	"context"
)

type Authenticator interface {
	Authenticate(ctx context.Context, username, password string) (scopes []string, success bool)
}

type authenticator struct {
	repo   repo.Repo
	hasher hasher.Hash
}

func NewAuthenticator(repo repo.Repo, hasher hasher.Hash) Authenticator {
	return &authenticator{
		repo:   repo,
		hasher: hasher,
	}
}

func (authenticatorService *authenticator) Authenticate(ctx context.Context, username, password string) ([]string, bool) {
	user, err := authenticatorService.repo.GetUser(ctx, username)
	if err != nil {
		return []string{}, false
	}
	passwordHashMatch := authenticatorService.hasher.ComparePassword(user.HashedPassword, password)
	if !passwordHashMatch {
		return []string{}, false
	}
	return user.Scopes, true
}
