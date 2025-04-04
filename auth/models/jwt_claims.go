package models

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	Username string   `json:"username"`
	Scopes   []string `json:"scopes"`
	jwt.RegisteredClaims
}
