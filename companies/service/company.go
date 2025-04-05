package service

import (
	"companies/consts"
	"companies/models"
	"companies/repo"
	"context"
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var CompaniesEventsKafkaTopic string = "companies-events"

type CompanyService interface {
	CreateCompany(ctx context.Context, companyInput models.CompanyInput) (models.CompanyOutput, error)
	PatchCompany(ctx context.Context, companyId uuid.UUID, updateCompanyInput models.UpdateCompanyInput) (models.CompanyOutput, error)
	GetCompany(ctx context.Context, companyId uuid.UUID) (models.CompanyOutput, error)
	DeleteCompany(ctx context.Context, companyId uuid.UUID) error
}

type companyService struct {
	repo          repo.CompanyRepo
	kafkaProducer *kafka.Producer
}

func NewCompanyService(repo repo.CompanyRepo, kafkaProducer *kafka.Producer) CompanyService {
	return &companyService{
		repo:          repo,
		kafkaProducer: kafkaProducer,
	}
}

func (service *companyService) CreateCompany(ctx context.Context, companyInput models.CompanyInput) (models.CompanyOutput, error) {
	company := models.Company{}
	company.ID = uuid.New()
	company.FromCompanyInput(companyInput)
	insertedId, err := service.repo.CreateCompany(ctx, company)
	if err != nil {
		return models.CompanyOutput{}, err
	}
	company.ID = insertedId
	output := models.CompanyOutput{
		ID: insertedId,
	}
	output.FromCompany(company)

	event := models.KafkaEvent{
		Type: models.KafkaEventTypeCompanyCreate,
		Data: output,
	}

	go service.publishEvent(event)

	return output, nil
}

func (service *companyService) PatchCompany(ctx context.Context, companyId uuid.UUID, updateCompanyInput models.UpdateCompanyInput) (models.CompanyOutput, error) {
	company, err := service.repo.PatchCompany(ctx, companyId, updateCompanyInput)
	if err != nil {
		return models.CompanyOutput{}, err
	}
	output := models.CompanyOutput{}
	output.FromCompany(company)

	event := models.KafkaEvent{
		Type: models.KafkaEventTypeCompanyPatch,
		Data: output,
	}

	go service.publishEvent(event)

	return output, nil
}

func (service *companyService) GetCompany(ctx context.Context, companyId uuid.UUID) (models.CompanyOutput, error) {
	company, err := service.repo.GetCompany(ctx, companyId)
	if err != nil {
		return models.CompanyOutput{}, err
	}
	companyOutput := models.CompanyOutput{}
	companyOutput.FromCompany(company)

	event := models.KafkaEvent{
		Type: models.KafkaEventTypeCompanyGet,
		Data: companyOutput,
	}

	go service.publishEvent(event)

	return companyOutput, nil
}

func (service *companyService) DeleteCompany(ctx context.Context, companyId uuid.UUID) error {
	err := service.repo.DeleteCompany(ctx, companyId)
	if err != nil {
		return err
	}

	event := models.KafkaEvent{
		Type: models.KafkaEventTypeCompanyDelete,
		Data: nil,
	}

	go service.publishEvent(event)

	return nil
}

func (service *companyService) publishEvent(event models.KafkaEvent) {
	if service.kafkaProducer == nil {
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

	err = service.kafkaProducer.Produce(
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
