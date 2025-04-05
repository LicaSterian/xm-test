package handlers

import (
	"auth/jwt"
	"auth/models"
	"auth/service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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
	err := c.ShouldBindJSON(&input)
	if err != nil {
		output.ErrorCode = ErrCodeInvalidInput
		c.JSON(http.StatusBadRequest, output)
		return
	}

	scopes, authenticated := handler.authenticatorService.Authenticate(ctx, input.Username, input.Password)
	if !authenticated {
		output.ErrorCode = ErrCodeAuthFailed
		c.JSON(http.StatusUnauthorized, output)
		return
	}

	token, err := handler.jwtGenerator.Generate(input.Username, scopes)
	if err != nil {
		output.ErrorCode = ErrCodeCouldNotCreateToken
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
	err := c.ShouldBindJSON(&input)
	if err != nil {
		output.ErrorCode = ErrCodeInvalidInput
		c.JSON(http.StatusBadRequest, output)
		return
	}

	err = handler.registratorService.Register(ctx, input.Username, input.Password, input.Scopes)
	if err != nil {
		output.ErrorCode = ErrCodeRegistrationFailed
		c.JSON(http.StatusInternalServerError, output)
		return
	}

	c.JSON(http.StatusOK, output)
}
