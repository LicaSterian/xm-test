package jwt

import (
	"auth/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTGenerator public interface
type JWTGenerator interface {
	Generate(username string, scopes []string) (string, error)
}

// jwtGenerator private struct that implements the JWTGenerator methods
type jwtGenerator struct {
	jwtKey []byte
	ttl    time.Duration
}

func NewJWTGenerator(key []byte, ttl time.Duration) JWTGenerator {
	return &jwtGenerator{
		jwtKey: key,
		ttl:    ttl,
	}
}

func (generator *jwtGenerator) Generate(username string, scopes []string) (string, error) {
	now := time.Now().UTC()
	expirationTime := now.Add(generator.ttl)

	claims := &models.JWTClaims{
		Username: username,
		Scopes:   scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(generator.jwtKey)
}
