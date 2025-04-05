package service

import (
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
	repo repo.CompanyRepo
}

func NewCompanyService(repo repo.CompanyRepo) CompanyService {
	return &companyService{
		repo: repo,
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
	return output, nil
}

func (service *companyService) PatchCompany(ctx context.Context, companyId uuid.UUID, updateCompanyInput models.UpdateCompanyInput) (models.CompanyOutput, error) {
	return models.CompanyOutput{}, nil
}

func (service *companyService) GetCompany(ctx context.Context, companyId uuid.UUID) (models.CompanyOutput, error) {
	company, err := service.repo.GetCompany(ctx, companyId)
	if err != nil {
		return models.CompanyOutput{}, err
	}
	companyOutput := models.CompanyOutput{}
	companyOutput.FromCompany(company)
	return companyOutput, nil
}

func (service *companyService) DeleteCompany(ctx context.Context, companyId uuid.UUID) error {
	return nil
}
