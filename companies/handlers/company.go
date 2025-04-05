package handlers

import (
	"companies/consts"
	"companies/models"
	"companies/service"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CompanyHandler interface {
	CreateCompany(c *gin.Context)
	PatchCompany(c *gin.Context)
	GetCompany(c *gin.Context)
	DeleteCompany(c *gin.Context)
}

type companyHandler struct {
	service service.CompanyService
}

func NewCompanyHandler(companyService service.CompanyService) CompanyHandler {
	return &companyHandler{
		service: companyService,
	}
}

func (handler *companyHandler) CreateCompany(c *gin.Context) {
	ctx := c.Request.Context()

	var companyInput models.CompanyInput
	err := c.ShouldBindJSON(&companyInput)
	if err != nil {
		errOutput := models.ErrorOutput{
			ErrorCode: ErrCodeInvalidInput,
		}
		err = errors.Join(ErrInvalidInput, err)
		log.Error().
			Err(err).
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Int(consts.LogKeyErrorCode, errOutput.ErrorCode).
			Int(consts.LogKeyStatusCode, http.StatusBadRequest).
			Msg("error while trying to bind JSON input")
		c.JSON(http.StatusBadRequest, errOutput)
		return
	}

	companyOutput, err := handler.service.CreateCompany(ctx, companyInput)
	if err != nil {
		errOutput := models.ErrorOutput{
			ErrorCode: ErrCodeInvalidInput,
		}
		err = errors.Join(ErrCouldNotCreateCompany, err)
		log.Error().
			Err(err).
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Int(consts.LogKeyErrorCode, errOutput.ErrorCode).
			Int(consts.LogKeyStatusCode, http.StatusInternalServerError).
			Msg("error while trying to create company")
		c.JSON(http.StatusInternalServerError, errOutput)
		return
	}

	log.Info().
		Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
		Int(consts.LogKeyStatusCode, http.StatusCreated).
		Msg("login successful")
	c.JSON(http.StatusCreated, companyOutput)
}

func (handler *companyHandler) PatchCompany(c *gin.Context) {

}

func (handler *companyHandler) GetCompany(c *gin.Context) {

}

func (handler *companyHandler) DeleteCompany(c *gin.Context) {

}
