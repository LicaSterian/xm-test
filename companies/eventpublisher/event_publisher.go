package eventpublisher

import (
	"companies/consts"
	"companies/models"
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
)

var CompaniesEventsKafkaTopic string = "companies-events"

type EventPublisher interface {
	PublishEvent(event models.KafkaEvent)
}

type eventPublisher struct {
	kafkaProducer *kafka.Producer
}

func NewEventPublisher(kafkaProducer *kafka.Producer) EventPublisher {
	return &eventPublisher{
		kafkaProducer: kafkaProducer,
	}
}

func (publisher *eventPublisher) PublishEvent(event models.KafkaEvent) {
	if publisher.kafkaProducer == nil {
		return
	}

	jsonEvent, err := json.Marshal(event)
	if err != nil {
		log.Error().
			Err(err).
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Str(consts.LogKeyKafkaEventType, string(event.Type)).
			Msg("error while marshalling JSON when trying to publish event")
		return
	}

	err = publisher.kafkaProducer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &CompaniesEventsKafkaTopic,
				Partition: kafka.PartitionAny,
			},
			Value: jsonEvent,
		},
		nil, // delivery channel
	)
	if err != nil {
		log.Error().
			Err(err).
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Str(consts.LogKeyKafkaEventType, string(event.Type)).
			Msg("error while publishing event")
		return
	}

	log.Info().
		Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
		Str(consts.LogKeyKafkaEventType, string(event.Type)).
		Msg("published event")
}
