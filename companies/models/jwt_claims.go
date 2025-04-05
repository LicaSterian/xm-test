package models

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Username string   `json:"username"`
	Scopes   []string `json:"scopes"`
	jwt.RegisteredClaims
}

func (c *JWTClaims) Valid() error {
	now := time.Now().UTC()
	if c.ExpiresAt != nil && c.ExpiresAt.Time.Before(now) {
		return errors.New("token is expired")
	}
	if c.NotBefore != nil && c.NotBefore.Time.After(now) {
		return errors.New("token not valid yet")
	}

	// TODO check issuer
	return nil
}
