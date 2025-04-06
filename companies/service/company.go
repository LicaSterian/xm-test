package service

import (
	"companies/eventpublisher"
	"companies/models"
	"companies/repo"
	"context"

	"github.com/google/uuid"
)

type CompanyService interface {
	CreateCompany(ctx context.Context, companyInput models.CompanyInput) (models.CompanyOutput, error)
	PatchCompany(ctx context.Context, companyId uuid.UUID, updateCompanyInput models.UpdateCompanyInput) (models.CompanyOutput, error)
	GetCompany(ctx context.Context, companyId uuid.UUID) (models.CompanyOutput, error)
	DeleteCompany(ctx context.Context, companyId uuid.UUID) error
}

type companyService struct {
	repo           repo.CompanyRepo
	eventPublisher eventpublisher.EventPublisher
}

func NewCompanyService(repo repo.CompanyRepo, eventPublisher eventpublisher.EventPublisher) CompanyService {
	return &companyService{
		repo:           repo,
		eventPublisher: eventPublisher,
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

	go service.eventPublisher.PublishEvent(event)

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

	go service.eventPublisher.PublishEvent(event)

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

	go service.eventPublisher.PublishEvent(event)

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

	go service.eventPublisher.PublishEvent(event)

	return nil
}
