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
	"github.com/rs/zerolog/log"
)

const (
	statusCodeLogKey = "status_code"
	errorCodeLogKey  = "error_code"
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
		err = errors.Join(ErrInvalidInput, err)
		output.ErrorCode = ErrCodeInvalidInput
		log.Error().
			Err(err).
			Int(errorCodeLogKey, output.ErrorCode).
			Int(statusCodeLogKey, http.StatusBadRequest).
			Msg("error while trying to bind JSON input")

		c.JSON(http.StatusBadRequest, output)
		return
	}

	scopes, authenticated := handler.authenticatorService.Authenticate(ctx, input.Username, input.Password)
	if !authenticated {
		output.ErrorCode = ErrCodeAuthFailed
		log.Error().
			Err(ErrAuthFailed).
			Int(errorCodeLogKey, output.ErrorCode).
			Int(statusCodeLogKey, http.StatusUnauthorized).
			Msg("error while trying to authenticate")
		c.JSON(http.StatusUnauthorized, output)
		return
	}

	token, err := handler.jwtGenerator.Generate(input.Username, scopes)
	if err != nil {
		err = errors.Join(ErrCouldNotGenerateToken, err)
		output.ErrorCode = ErrCodeCouldNotGenerateToken
		log.Error().
			Err(err).
			Int(errorCodeLogKey, output.ErrorCode).
			Int(statusCodeLogKey, http.StatusInternalServerError).
			Msg("error while trying to generate token")
		c.JSON(http.StatusInternalServerError, output)
		return
	}

	output.Token = token
	log.Info().
		Int(statusCodeLogKey, http.StatusOK).
		Msg("login successful")
	c.JSON(http.StatusOK, output)
}

func (handler *authHandler) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var input models.RegisterInput
	output := models.RegisterOutput{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err = errors.Join(ErrInvalidInput, err)
		output.ErrorCode = ErrCodeInvalidInput
		log.Error().
			Err(err).
			Int(errorCodeLogKey, output.ErrorCode).
			Int(statusCodeLogKey, http.StatusBadRequest).
			Msg("error while trying to bind JSON input")
		c.JSON(http.StatusBadRequest, output)
		return
	}

	err = handler.registratorService.Register(ctx, input.Username, input.Password, input.Scopes)
	if err != nil {
		err = errors.Join(ErrRegistrationFailed, err)
		output.ErrorCode = ErrCodeRegistrationFailed
		log.Error().
			Err(err).
			Int(errorCodeLogKey, output.ErrorCode).
			Int(statusCodeLogKey, http.StatusInternalServerError).
			Msg("error while trying to register")
		c.JSON(http.StatusInternalServerError, output)
		return
	}
	log.Info().
		Int(statusCodeLogKey, http.StatusOK).
		Msg("register successful")
	c.JSON(http.StatusOK, output)
}
