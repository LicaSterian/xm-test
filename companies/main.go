package main

import (
	"companies/consts"
	"companies/eventpublisher"
	"companies/handlers"
	"companies/middleware"
	"companies/repo"
	"companies/service"
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

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		err := errors.New("MONGO_URI env var not set")
		log.Error().
			Err(err).
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Msg("make sure to set the MONGO_URI env var")
		return
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		err := errors.New("JWT_SECRET_KEY env var not set")
		log.Error().
			Err(err).
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Msg("make sure to set the JWT_SECRET_KEY env var")
		return
	}

	kafkaServers := os.Getenv("KAFKA_SERVERS")
	if kafkaServers == "" {
		err := errors.New("KAFKA_SERVERS env var not set. ex: KAFKA_SERVERS=localhost:9092,example.com:9092")
		log.Error().
			Err(err).
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Msg("make sure to set the KAFKA_SERVERS env var")
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
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Msg("failed to connect to MongoDB")
		return
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error().
			Err(err).
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Msg("failed to ping MongoDB")
		return
	}

	log.Info().
		Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
		Msg("connected to MongoDB")

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaServers,
		"client.id":         "companies-service",
		"acks":              "all"})

	if err != nil {
		log.Error().
			Err(err).
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Msg("failed to create Kafka producer")
		return
	}

	log.Info().
		Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
		Msg("created Kafka producer")

	// Setup the service

	companyRepo := repo.NewMongoCompanyRepo(client)
	eventPublisher := eventpublisher.NewEventPublisher(producer)
	companyService := service.NewCompanyService(companyRepo, eventPublisher)
	companyHandler := handlers.NewCompanyHandler(companyService)

	// setup gin engine
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.TimeoutMiddleware(5 * time.Second))

	v1Group := engine.Group("/v1", middleware.ValidateJWTToken([]byte(jwtSecretKey)))

	v1Group.POST("/company", companyHandler.CreateCompany)
	v1Group.PATCH("/company/:id", companyHandler.PatchCompany)
	v1Group.GET("/company/:id", companyHandler.GetCompany)
	v1Group.DELETE("/company/:id", companyHandler.DeleteCompany)

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// run the service on port 8080
	go func() {
		port := "8080"
		err := engine.Run(":" + port)
		if err != nil {
			log.Error().
				Err(err).
				Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
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
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Msg("error while disconnecting from MongoDB")
	}

	log.Info().
		Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
		Msg("successfully disconnected from MongoDB")

	// close kafka connection
	producer.Close()

	log.Info().
		Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
		Msg("closed the Kafka producer")
}
