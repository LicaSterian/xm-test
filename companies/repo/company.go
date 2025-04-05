package repo

import (
	"companies/models"
	"context"

	"github.com/google/uuid"
)

type CompanyRepo interface {
	CreateCompany(ctx context.Context, company models.Company) (uuid.UUID, error)
	PatchCompany(ctx context.Context, companyId uuid.UUID, company models.UpdateCompanyInput) (models.Company, error)
	GetCompany(ctx context.Context, companyId uuid.UUID) (models.Company, error)
	DeleteCompany(ctx context.Context, companyId uuid.UUID) error
}
