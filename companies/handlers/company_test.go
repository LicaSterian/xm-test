package handlers

import (
	"bytes"
	"companies/mocks"
	"companies/models"
	"companies/repo"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCreateCompany(t *testing.T) {
	companyId := uuid.New()

	testCases := []struct {
		name                 string
		requestBody          string
		companyOutput        models.CompanyOutput
		expectedStatusCode   int
		expectedResponseBody string
		stubMocks            func(s *mocks.CompanyService, companyOutput models.CompanyOutput)
	}{
		{
			name: "success test case",
			requestBody: `{
				"name": "company-name",
				"description": "company-description",
				"number_of_employees": 100,
				"registered": true,
				"type": "Corporations"
			}`,
			companyOutput: models.CompanyOutput{
				ID:                companyId,
				Name:              "company-name",
				Description:       "company-description",
				NumberOfEmployees: 100,
				Registered:        true,
				Type:              "Corporations",
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponseBody: fmt.Sprintf(`{
				"id": "%s",
				"name":"company-name",
				"description":"company-description",
				"number_of_employees": 100,
				"registered": true,
				"type": "Corporations"
			}`, companyId),
			stubMocks: func(s *mocks.CompanyService, companyOutput models.CompanyOutput) {
				s.On("CreateCompany", mock.Anything, mock.AnythingOfType("models.CompanyInput")).
					Return(companyOutput, nil)
			},
		},
		{
			name: "name field is required",
			requestBody: `{
				"description": "company-description",
				"number_of_employees": 100,
				"registered": true,
				"type": "Corporations"
			}`,
			companyOutput:      models.CompanyOutput{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: fmt.Sprintf(`{
				"error_code": %d
			}`, ErrCodeInvalidInput),
			stubMocks: func(s *mocks.CompanyService, companyOutput models.CompanyOutput) {

			},
		},
		// TODO test case for each validation field
		{
			name: "xss content in description",
			requestBody: `{
				"name": "company-name",
				"description": "<a onblur=\"alert('secret')\" href=\"http://www.google.com\">Google</a>",
				"number_of_employees": 100,
				"registered": true,
				"type": "Corporations"
			}`,
			companyOutput:      models.CompanyOutput{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: fmt.Sprintf(`{
				"error_code": %d
			}`, ErrCodeInvalidInput),
			stubMocks: func(s *mocks.CompanyService, companyOutput models.CompanyOutput) {

			},
		},
		{
			name: "service returns an 500 error",
			requestBody: `{
				"name": "company-name",
				"description": "company-description",
				"number_of_employees": 100,
				"registered": true,
				"type": "Corporations"
			}`,
			companyOutput:      models.CompanyOutput{},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponseBody: fmt.Sprintf(`{
				"error_code": %d
			}`, ErrCodeCouldNotCreateCompany),
			stubMocks: func(s *mocks.CompanyService, companyOutput models.CompanyOutput) {
				s.On("CreateCompany", mock.Anything, mock.AnythingOfType("models.CompanyInput")).
					Return(models.CompanyOutput{}, assert.AnError)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			s := new(mocks.CompanyService)

			handler := NewCompanyHandler(s)

			testCase.stubMocks(s, testCase.companyOutput)

			// Set Gin to test mode
			gin.SetMode(gin.TestMode)

			// Create a test context and response recorder
			router := gin.Default()
			router.POST("/v1/company", handler.CreateCompany)

			buf := bytes.NewBuffer([]byte(testCase.requestBody))

			req, _ := http.NewRequest(http.MethodPost, "/v1/company", buf)
			req.Header.Set("content-type", "application/json")
			rr := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(rr, req)

			// Assertions
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, rr.Body.String())
		})
	}
}

func TestPatchCompany(t *testing.T) {
	companyId := uuid.New()

	testCases := []struct {
		name                 string
		companyId            string
		requestBody          string
		companyOutput        models.CompanyOutput
		expectedStatusCode   int
		expectedResponseBody string
		stubMocks            func(s *mocks.CompanyService, companyOutput models.CompanyOutput)
	}{
		{
			name:      "success test case",
			companyId: companyId.String(),
			requestBody: `{
				"name": "company-name",
				"description": "company-description-updated",
				"number_of_employees": 100,
				"registered": false,
				"type": "NonProfit"
			}`,
			companyOutput: models.CompanyOutput{
				ID:                companyId,
				Name:              "company-name",
				Description:       "company-description-updated",
				NumberOfEmployees: 100,
				Registered:        false,
				Type:              "NonProfit",
			},
			expectedStatusCode: http.StatusAccepted,
			expectedResponseBody: fmt.Sprintf(`{
				"id": "%s",
				"name":"company-name",
				"description":"company-description-updated",
				"number_of_employees": 100,
				"registered": false,
				"type": "NonProfit"
			}`, companyId.String()),
			stubMocks: func(s *mocks.CompanyService, companyOutput models.CompanyOutput) {
				s.On("PatchCompany", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("models.UpdateCompanyInput")).
					Return(companyOutput, nil)
			},
		},
		{
			name:               "invalid companyId",
			companyId:          companyId.String() + "abc",
			requestBody:        ``,
			companyOutput:      models.CompanyOutput{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: fmt.Sprintf(`{
				"error_code": %d
			}`, ErrCodeInvalidId),
			stubMocks: func(s *mocks.CompanyService, companyOutput models.CompanyOutput) {

			},
		},
		{
			name:      "name field longer than 15 chars",
			companyId: companyId.String(),
			requestBody: `{
				"name": "company-name-longer-than-15-chars"
			}`,
			companyOutput:      models.CompanyOutput{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: fmt.Sprintf(`{
				"error_code": %d
			}`, ErrCodeInvalidInput),
			stubMocks: func(s *mocks.CompanyService, companyOutput models.CompanyOutput) {

			},
		},
		{
			name:      "xss content in description",
			companyId: companyId.String(),
			requestBody: `{
				"description": "<a onblur=\"alert('secret')\" href=\"http://www.google.com\">Google</a>"
			}`,
			companyOutput:      models.CompanyOutput{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: fmt.Sprintf(`{
				"error_code": %d
			}`, ErrCodeInvalidInput),
			stubMocks: func(s *mocks.CompanyService, companyOutput models.CompanyOutput) {

			},
		},
		{
			name:      "test case 404",
			companyId: companyId.String(),
			requestBody: `{
				"name": "company-name"
			}`,
			companyOutput:      models.CompanyOutput{},
			expectedStatusCode: http.StatusNotFound,
			expectedResponseBody: fmt.Sprintf(`{
				"error_code": %d
			}`, ErrCodePatchCompany),
			stubMocks: func(s *mocks.CompanyService, companyOutput models.CompanyOutput) {
				s.On("PatchCompany", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("models.UpdateCompanyInput")).
					Return(models.CompanyOutput{}, errors.Join(assert.AnError, mongo.ErrNoDocuments))
			},
		},
		{
			name:      "test case 500",
			companyId: companyId.String(),
			requestBody: `{
				"name": "company-name"
			}`,
			companyOutput:      models.CompanyOutput{},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponseBody: fmt.Sprintf(`{
				"error_code": %d
			}`, ErrCodePatchCompany),
			stubMocks: func(s *mocks.CompanyService, companyOutput models.CompanyOutput) {
				s.On("PatchCompany", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("models.UpdateCompanyInput")).
					Return(models.CompanyOutput{}, assert.AnError)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			s := new(mocks.CompanyService)

			handler := NewCompanyHandler(s)

			testCase.stubMocks(s, testCase.companyOutput)

			// Set Gin to test mode
			gin.SetMode(gin.TestMode)

			// Create a test context and response recorder
			router := gin.Default()
			router.PATCH("/v1/company/:id", handler.PatchCompany)

			buf := bytes.NewBuffer([]byte(testCase.requestBody))

			url := fmt.Sprintf("/v1/company/%s", testCase.companyId)
			req, _ := http.NewRequest(http.MethodPatch, url, buf)
			req.Header.Set("content-type", "application/json")
			rr := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(rr, req)

			// Assertions
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, rr.Body.String())
		})
	}
}

func TestGetCompany(t *testing.T) {
	companyId := uuid.New()

	testCases := []struct {
		name                 string
		companyId            string
		companyOutput        models.CompanyOutput
		expectedStatusCode   int
		expectedResponseBody string
		stubMocks            func(s *mocks.CompanyService, companyOutput models.CompanyOutput)
	}{
		{
			name:      "success test case",
			companyId: companyId.String(),
			companyOutput: models.CompanyOutput{
				ID:                companyId,
				Name:              "company-name",
				Description:       "company-description",
				NumberOfEmployees: 100,
				Registered:        true,
				Type:              "Corporations",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: fmt.Sprintf(`{
				"id": "%s",
				"name":"company-name",
				"description":"company-description",
				"number_of_employees": 100,
				"registered": true,
				"type": "Corporations"
			}`, companyId.String()),
			stubMocks: func(s *mocks.CompanyService, companyOutput models.CompanyOutput) {
				s.On("GetCompany", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(companyOutput, nil)
			},
		},
		{
			name:               "invalid companyId",
			companyId:          "abc",
			companyOutput:      models.CompanyOutput{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: fmt.Sprintf(`{
				"error_code": %d
			}`, ErrCodeInvalidId),
			stubMocks: func(s *mocks.CompanyService, companyOutput models.CompanyOutput) {

			},
		},
		{
			name:               "test case 404",
			companyId:          companyId.String(),
			companyOutput:      models.CompanyOutput{},
			expectedStatusCode: http.StatusNotFound,
			expectedResponseBody: fmt.Sprintf(`{
				"error_code": %d
			}`, ErrCodeGetCompany),
			stubMocks: func(s *mocks.CompanyService, companyOutput models.CompanyOutput) {
				s.On("GetCompany", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(models.CompanyOutput{}, errors.Join(assert.AnError, mongo.ErrNoDocuments))
			},
		},
		{
			name:               "test case 500",
			companyId:          companyId.String(),
			companyOutput:      models.CompanyOutput{},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponseBody: fmt.Sprintf(`{
				"error_code": %d
			}`, ErrCodeGetCompany),
			stubMocks: func(s *mocks.CompanyService, companyOutput models.CompanyOutput) {
				s.On("GetCompany", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(models.CompanyOutput{}, assert.AnError)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			s := new(mocks.CompanyService)

			handler := NewCompanyHandler(s)

			testCase.stubMocks(s, testCase.companyOutput)

			// Set Gin to test mode
			gin.SetMode(gin.TestMode)

			// Create a test context and response recorder
			router := gin.Default()
			router.GET("/v1/company/:id", handler.GetCompany)

			url := fmt.Sprintf("/v1/company/%s", testCase.companyId)
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("content-type", "application/json")
			rr := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(rr, req)

			// Assertions
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, rr.Body.String())
		})
	}
}

func TestDeleteCompany(t *testing.T) {
	companyId := uuid.New()

	testCases := []struct {
		name               string
		companyId          string
		expectedStatusCode int
		stubMocks          func(s *mocks.CompanyService)
	}{
		{
			name:               "success test case",
			companyId:          companyId.String(),
			expectedStatusCode: http.StatusNoContent,
			stubMocks: func(s *mocks.CompanyService) {
				s.On("DeleteCompany", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(nil)
			},
		},
		{
			name:               "invalid companyId",
			companyId:          "abc",
			expectedStatusCode: http.StatusBadRequest,
			stubMocks: func(s *mocks.CompanyService) {

			},
		},
		{
			name:               "test case 404",
			companyId:          companyId.String(),
			expectedStatusCode: http.StatusNotFound,
			stubMocks: func(s *mocks.CompanyService) {
				s.On("DeleteCompany", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(errors.Join(assert.AnError, repo.ErrDocumentNotFound))
			},
		},
		{
			name:               "test case 500",
			companyId:          companyId.String(),
			expectedStatusCode: http.StatusInternalServerError,
			stubMocks: func(s *mocks.CompanyService) {
				s.On("DeleteCompany", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(assert.AnError)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			s := new(mocks.CompanyService)

			handler := NewCompanyHandler(s)

			testCase.stubMocks(s)

			// Set Gin to test mode
			gin.SetMode(gin.TestMode)

			// Create a test context and response recorder
			router := gin.Default()
			router.DELETE("/v1/company/:id", handler.DeleteCompany)

			url := fmt.Sprintf("/v1/company/%s", testCase.companyId)
			req, _ := http.NewRequest(http.MethodDelete, url, nil)
			req.Header.Set("content-type", "application/json")
			rr := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(rr, req)

			// Assertions
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
		})
	}
}
