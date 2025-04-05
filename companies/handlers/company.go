package handlers

import (
	"companies/consts"
	"companies/models"
	"companies/service"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
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
	ctx := c.Request.Context()

	companyIdParam := c.Param("id")
	companyId, err := uuid.Parse(companyIdParam)
	if err != nil {
		errOutput := models.ErrorOutput{
			ErrorCode: ErrCodeInvalidId,
		}
		err = errors.Join(ErrInvalidInput, err)
		log.Error().
			Err(err).
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Int(consts.LogKeyErrorCode, errOutput.ErrorCode).
			Int(consts.LogKeyStatusCode, http.StatusBadRequest).
			Str(consts.LogKeyCompanyId, companyIdParam).
			Msg("error while trying to parse companyId")
		c.JSON(http.StatusBadRequest, errOutput)
		return
	}

	var updateCompanyInput models.UpdateCompanyInput
	err = c.ShouldBindJSON(&updateCompanyInput)
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
			Str(consts.LogKeyCompanyId, companyId.String()).
			Msg("error while trying to bind JSON input")
		c.JSON(http.StatusBadRequest, errOutput)
		return
	}

	companyOutput, err := handler.service.PatchCompany(ctx, companyId, updateCompanyInput)
	if err != nil {
		errOutput := models.ErrorOutput{
			ErrorCode: ErrCodePatchCompany,
		}
		err = errors.Join(ErrPatchCompany, err)
		log.Error().
			Err(err).
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Int(consts.LogKeyErrorCode, errOutput.ErrorCode).
			Int(consts.LogKeyStatusCode, http.StatusBadRequest).
			Str(consts.LogKeyCompanyId, companyId.String()).
			Msg("error while trying to PATCH company")
		c.JSON(http.StatusBadRequest, errOutput)
		return
	}

	log.Info().
		Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
		Int(consts.LogKeyStatusCode, http.StatusOK).
		Str(consts.LogKeyCompanyId, companyId.String()).
		Msg("patch company executed successfully")
	c.JSON(http.StatusAccepted, companyOutput)
}

func (handler *companyHandler) GetCompany(c *gin.Context) {
	ctx := c.Request.Context()

	companyIdParam := c.Param("id")
	companyId, err := uuid.Parse(companyIdParam)
	if err != nil {
		errOutput := models.ErrorOutput{
			ErrorCode: ErrCodeInvalidId,
		}
		err = errors.Join(ErrInvalidId, err)
		log.Error().
			Err(err).
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Int(consts.LogKeyErrorCode, errOutput.ErrorCode).
			Int(consts.LogKeyStatusCode, http.StatusBadRequest).
			Str(consts.LogKeyCompanyId, companyIdParam).
			Msg("error while trying to parse companyId")
		c.JSON(http.StatusBadRequest, errOutput)
		return
	}

	companyOutput, err := handler.service.GetCompany(ctx, companyId)
	if err != nil {
		errOutput := models.ErrorOutput{
			ErrorCode: ErrCodeGetCompany,
		}
		err = errors.Join(ErrGetCompany, err)
		statusCode := http.StatusInternalServerError
		if errors.Is(err, mongo.ErrNoDocuments) {
			statusCode = http.StatusNotFound
		}
		log.Error().
			Err(err).
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Int(consts.LogKeyErrorCode, errOutput.ErrorCode).
			Int(consts.LogKeyStatusCode, statusCode).
			Str(consts.LogKeyCompanyId, companyId.String()).
			Msg("error while trying to get company")
		c.JSON(statusCode, errOutput)
		return
	}

	log.Info().
		Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
		Int(consts.LogKeyStatusCode, http.StatusOK).
		Str(consts.LogKeyCompanyId, companyId.String()).
		Msg("get company executed successfully")
	c.JSON(http.StatusOK, companyOutput)
}

func (handler *companyHandler) DeleteCompany(c *gin.Context) {

}
