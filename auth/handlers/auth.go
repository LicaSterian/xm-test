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

const errKeyInvalidInput = "invalid input"
const errKeyAuthFailed = "authentication failed"
const errKeyCouldNotCreateToken = "could not create token"
const errKeyRegistrationFailed = "registration failed"

var ErrInvalidInput = errors.New(errKeyInvalidInput)
var ErrAuthFailed = errors.New(errKeyAuthFailed)
var ErrCouldNotCreateToken = errors.New(errKeyCouldNotCreateToken)

var errorCodes = map[string]int{
	errKeyInvalidInput:        1,
	errKeyAuthFailed:          2,
	errKeyCouldNotCreateToken: 3,
}

type AuthHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type authHandler struct {
	authenticatorService service.Authenticator
	registratorService   service.Registrator
	jwtGenerator         jwt.JWTGenerator
}

func NewAuthHandler(
	authenticatorService service.Authenticator,
	registratorService service.Registrator,
	jwtGenerator jwt.JWTGenerator,
) AuthHandler {
	return &authHandler{
		authenticatorService: authenticatorService,
		registratorService:   registratorService,
		jwtGenerator:         jwtGenerator,
	}
}

func (handler *authHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var input models.LoginInput
	output := models.LoginOutput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		output.ErrorCode = errorCodes[errKeyInvalidInput]
		c.JSON(http.StatusBadRequest, output)
		return
	}

	scopes, authenticated := handler.authenticatorService.Authenticate(ctx, input.Username, input.Password)
	if !authenticated {
		output.ErrorCode = errorCodes[errKeyAuthFailed]
		c.JSON(http.StatusUnauthorized, output)
		return
	}

	token, err := handler.jwtGenerator.Generate(input.Username, scopes)
	if err != nil {
		output.ErrorCode = errorCodes[errKeyCouldNotCreateToken]
		c.JSON(http.StatusInternalServerError, output)
		return
	}

	output.Token = token
	c.JSON(http.StatusOK, output)
}

func (handler *authHandler) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var input models.RegisterInput
	output := models.RegisterOutput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		output.ErrorCode = errorCodes[errKeyInvalidInput]
		c.JSON(http.StatusBadRequest, output)
		return
	}

	err := handler.registratorService.Register(ctx, input.Username, input.Password, input.Scopes)
	if err != nil {
		output.ErrorCode = errorCodes[errKeyRegistrationFailed]
		c.JSON(http.StatusInternalServerError, output)
		return
	}

	c.JSON(http.StatusOK, output)
}
