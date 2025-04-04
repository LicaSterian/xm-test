package main

import (
	"auth/handlers"
	"auth/hasher"
	"auth/jwt"
	"auth/repo"
	"auth/service"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI env var not set")
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatal("JWT_SECRET_KEY env var not set")
	}

	// Set up a connection to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB")

	repo := repo.NewMongoRepo(client)
	hasher := hasher.NewHasher()

	authenticatorService := service.NewAuthenticator(repo, hasher)
	jwtGenerator := jwt.NewJWTGenerator([]byte(jwtSecretKey), time.Hour)

	authHandler := handlers.NewAuthHandler(authenticatorService, jwtGenerator)

	// setup gin engine
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.POST("/login", authHandler.Login)

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Simulate some work
	go func() {
		err := engine.Run(":80")
		if err != nil {
			log.Panic("error running the server on port 80:", err.Error())
		}
	}()

	// Wait until context is canceled
	<-ctx.Done()

	// close mongodb connection
	disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer disconnectCancel()

	err = client.Disconnect(disconnectCtx)
	if err != nil {
		log.Panic("error while disconnecting from mongodb", err.Error())
	}
}
