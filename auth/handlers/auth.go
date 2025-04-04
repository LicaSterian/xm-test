package handlers

import (
	"auth/jwt"
	"auth/models"
	"auth/service"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const errKeyAuthFailed = "authentication failed"
const errKeyCouldNotCreateToken = "could not create token"

var ErrAuthFailed = errors.New(errKeyAuthFailed)
var ErrCouldNotCreateToken = errors.New(errKeyCouldNotCreateToken)

var errorCodes = map[string]int{
	errKeyAuthFailed:          1,
	errKeyCouldNotCreateToken: 2,
}

type AuthHandler interface {
	Login(c *gin.Context)
}

type authHandler struct {
	authenticatorService service.Authenticator
	jwtGenerator         jwt.JWTGenerator
}

func NewAuthHandler(authenticatorService service.Authenticator, jwtGenerator jwt.JWTGenerator) AuthHandler {
	return &authHandler{
		authenticatorService: authenticatorService,
		jwtGenerator:         jwtGenerator,
	}
}

func (handler *authHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	output := models.LoginOutput{}

	if !handler.authenticatorService.Authenticate(ctx, input.Username, input.Password) {
		output.ErrorCode = errorCodes[errKeyAuthFailed]
		c.JSON(http.StatusUnauthorized, output)
		return
	}

	token, err := handler.jwtGenerator.Generate(input.Username)
	if err != nil {
		output.ErrorCode = errorCodes[errKeyCouldNotCreateToken]
		c.JSON(http.StatusInternalServerError, output)
		return
	}

	output.Token = token
	c.JSON(http.StatusOK, output)
}
