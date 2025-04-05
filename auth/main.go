package main

import (
	"auth/handlers"
	"auth/hasher"
	"auth/jwt"
	"auth/repo"
	"auth/service"
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		err := errors.New("MONGO_URI env var not set")
		log.Error().
			Err(err).Msg("make sure to set the MONGO_URI env var")
		return
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		err := errors.New("JWT_SECRET_KEY env var not set")
		log.Error().
			Err(err).
			Msg("make sure to set the JWT_SECRET_KEY env var")
		return
	}

	// Set up a connection to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to connect to MongoDB")
		return
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to ping MongoDB")
		return
	}

	log.Info().
		Msg("connected to MongoDB")

	repo := repo.NewMongoRepo(client)
	hasher := hasher.NewHasher()

	authenticatorService := service.NewAuthenticator(repo, hasher)
	registratorService := service.NewRegistrator(repo, hasher)
	jwtGenerator := jwt.NewJWTGenerator([]byte(jwtSecretKey), time.Hour)

	authHandler := handlers.NewAuthHandler(authenticatorService, registratorService, jwtGenerator)

	// setup gin engine
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.POST("/login", authHandler.Login)
	engine.POST("register", authHandler.Register)

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Simulate some work
	go func() {
		port := "80"
		err := engine.Run(":" + port)
		if err != nil {
			log.Error().
				Err(err).
				Msgf("error running the server on port %s", port)
		}
	}()

	// Wait until context is canceled
	<-ctx.Done()

	// close mongodb connection
	disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer disconnectCancel()

	err = client.Disconnect(disconnectCtx)
	if err != nil {
		log.Error().
			Err(err).
			Msg("error while disconnecting from MongoDB")
	}

	log.Info().
		Msg("successfully disconnected from MongoDB")
}
