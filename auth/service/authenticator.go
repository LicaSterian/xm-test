package service

import (
	"auth/hasher"
	"auth/repo"
	"context"
)

type Authenticator interface {
	Authenticate(ctx context.Context, username, password string) bool
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

func (authenticatorService *authenticator) Authenticate(ctx context.Context, username, password string) bool {
	hashedPassword, err := authenticatorService.repo.GetHashedPasswordByUsername(ctx, username)
	if err != nil {
		return false
	}
	return authenticatorService.hasher.ComparePassword(hashedPassword, password)
}
