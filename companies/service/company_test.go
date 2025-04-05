package service

import (
	"companies/mocks"
	"companies/models"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCompany(t *testing.T) {
	numberOfEmployees := 10
	registered := true

	testCases := []struct {
		name         string
		companyInput models.CompanyInput
		company      models.Company
		stubMock     func(r *mocks.CompanyRepo, company models.Company)
		validate     func(company models.Company, companyOutput models.CompanyOutput, err error)
	}{
		{
			name: "success test case",
			companyInput: models.CompanyInput{
				Name:              "company-name",
				Description:       "company-description",
				NumberOfEmployees: &numberOfEmployees,
				Registered:        &registered,
				Type:              "Corporations",
			},
			company: models.Company{
				ID:                uuid.New(),
				Name:              "company-name",
				Description:       "company-description",
				NumberOfEmployees: 10,
				Registered:        true,
				Type:              "Corporations",
			},
			stubMock: func(r *mocks.CompanyRepo, company models.Company) {
				r.On("CreateCompany", mock.Anything, mock.AnythingOfType("models.Company")).
					Return(company.ID, nil)
			},
			validate: func(company models.Company, companyOutput models.CompanyOutput, err error) {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, companyOutput.ID)

				expectedCompanyOutput := models.CompanyOutput{}
				expectedCompanyOutput.FromCompany(company)

				assert.Equal(t, expectedCompanyOutput, companyOutput)
			},
		},
		{
			name: "repo error",
			companyInput: models.CompanyInput{
				Name:              "company-name",
				Description:       "company-description",
				NumberOfEmployees: &numberOfEmployees,
				Registered:        &registered,
				Type:              "Corporations",
			},
			company: models.Company{
				ID:                uuid.New(),
				Name:              "company-name",
				Description:       "company-description",
				NumberOfEmployees: 10,
				Registered:        true,
				Type:              "Corporations",
			},
			stubMock: func(r *mocks.CompanyRepo, company models.Company) {
				r.On("CreateCompany", mock.Anything, mock.AnythingOfType("models.Company")).
					Return(uuid.Nil, assert.AnError)
			},
			validate: func(company models.Company, companyOutput models.CompanyOutput, err error) {
				assert.Error(t, err)
				assert.Equal(t, uuid.Nil, companyOutput.ID)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := new(mocks.CompanyRepo)

			companyService := NewCompanyService(r, nil)

			testCase.stubMock(r, testCase.company)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			companyOutput, err := companyService.CreateCompany(ctx, testCase.companyInput)
			testCase.validate(testCase.company, companyOutput, err)
		})
	}
}

func TestPatchCompany(t *testing.T) {
	companyName := "company-name"
	companyDescription := "company-description"
	numberOfEmployees := 10
	registered := true
	companyType := "Corporations"

	testCases := []struct {
		name               string
		companyId          uuid.UUID
		updateCompanyInput models.UpdateCompanyInput
		company            models.Company
		stubMock           func(r *mocks.CompanyRepo, company models.Company)
		validate           func(company models.Company, companyOutput models.CompanyOutput, err error)
	}{
		{
			name:      "success test case",
			companyId: uuid.New(),
			updateCompanyInput: models.UpdateCompanyInput{
				Name:              &companyName,
				Description:       &companyDescription,
				NumberOfEmployees: &numberOfEmployees,
				Registered:        &registered,
				Type:              &companyType,
			},
			company: models.Company{
				ID:                uuid.New(),
				Name:              companyName,
				Description:       companyDescription,
				NumberOfEmployees: 10,
				Registered:        true,
				Type:              companyType,
			},
			stubMock: func(r *mocks.CompanyRepo, company models.Company) {
				r.On("PatchCompany", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("models.UpdateCompanyInput")).
					Return(company, nil)
			},
			validate: func(company models.Company, companyOutput models.CompanyOutput, err error) {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, companyOutput.ID)

				expectedCompanyOutput := models.CompanyOutput{}
				expectedCompanyOutput.FromCompany(company)

				assert.Equal(t, expectedCompanyOutput, companyOutput)
			},
		},
		{
			name:      "repo returned an error",
			companyId: uuid.New(),
			updateCompanyInput: models.UpdateCompanyInput{
				Name:              &companyName,
				Description:       &companyDescription,
				NumberOfEmployees: &numberOfEmployees,
				Registered:        &registered,
				Type:              &companyType,
			},
			company: models.Company{
				ID:                uuid.New(),
				Name:              companyName,
				Description:       companyDescription,
				NumberOfEmployees: numberOfEmployees,
				Registered:        registered,
				Type:              companyType,
			},
			stubMock: func(r *mocks.CompanyRepo, company models.Company) {
				r.On("PatchCompany", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("models.UpdateCompanyInput")).
					Return(models.Company{}, assert.AnError)
			},
			validate: func(company models.Company, companyOutput models.CompanyOutput, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := new(mocks.CompanyRepo)

			companyService := NewCompanyService(r, nil)

			testCase.stubMock(r, testCase.company)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			companyOutput, err := companyService.PatchCompany(ctx, testCase.companyId, testCase.updateCompanyInput)
			testCase.validate(testCase.company, companyOutput, err)
		})
	}
}

func TestGetCompany(t *testing.T) {
	companyName := "company-name"
	companyDescription := "company-description"
	numberOfEmployees := 10
	registered := true
	companyType := "Corporations"

	testCases := []struct {
		name      string
		companyId uuid.UUID
		company   models.Company
		stubMock  func(r *mocks.CompanyRepo, company models.Company)
		validate  func(company models.Company, companyOutput models.CompanyOutput, err error)
	}{
		{
			name:      "success test case",
			companyId: uuid.New(),
			company: models.Company{
				ID:                uuid.New(),
				Name:              companyName,
				Description:       companyDescription,
				NumberOfEmployees: numberOfEmployees,
				Registered:        registered,
				Type:              companyType,
			},
			stubMock: func(r *mocks.CompanyRepo, company models.Company) {
				r.On("GetCompany", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(company, nil)
			},
			validate: func(company models.Company, companyOutput models.CompanyOutput, err error) {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, companyOutput.ID)

				expectedCompanyOutput := models.CompanyOutput{}
				expectedCompanyOutput.FromCompany(company)

				assert.Equal(t, expectedCompanyOutput, companyOutput)
			},
		},
		{
			name:      "repo returned an error",
			companyId: uuid.New(),
			company: models.Company{
				ID:                uuid.New(),
				Name:              companyName,
				Description:       companyDescription,
				NumberOfEmployees: numberOfEmployees,
				Registered:        registered,
				Type:              companyType,
			},
			stubMock: func(r *mocks.CompanyRepo, company models.Company) {
				r.On("GetCompany", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(models.Company{}, assert.AnError)
			},
			validate: func(company models.Company, companyOutput models.CompanyOutput, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := new(mocks.CompanyRepo)

			companyService := NewCompanyService(r, nil)

			testCase.stubMock(r, testCase.company)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			companyOutput, err := companyService.GetCompany(ctx, testCase.companyId)
			testCase.validate(testCase.company, companyOutput, err)
		})
	}
}

func TestDeleteCompany(t *testing.T) {
	testCases := []struct {
		name      string
		companyId uuid.UUID
		stubMock  func(r *mocks.CompanyRepo)
		validate  func(err error)
	}{
		{
			name:      "success test case",
			companyId: uuid.New(),
			stubMock: func(r *mocks.CompanyRepo) {
				r.On("DeleteCompany", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(nil)
			},
			validate: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:      "repo returned an error",
			companyId: uuid.New(),
			stubMock: func(r *mocks.CompanyRepo) {
				r.On("DeleteCompany", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(assert.AnError)
			},
			validate: func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := new(mocks.CompanyRepo)

			companyService := NewCompanyService(r, nil)

			testCase.stubMock(r)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := companyService.DeleteCompany(ctx, testCase.companyId)
			testCase.validate(err)
		})
	}
}
