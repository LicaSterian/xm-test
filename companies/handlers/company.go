package handlers

import (
	"companies/consts"
	"companies/models"
	"companies/repo"
	"companies/service"
	"companies/xss"
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

	// Name field only has 15 chars, it is not long enough to contain XSS content
	err = xss.CheckForXSS(companyInput.Description)
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
			Msg("error while checking for XSS content")
		c.JSON(http.StatusBadRequest, errOutput)
		return
	}

	companyOutput, err := handler.service.CreateCompany(ctx, companyInput)
	if err != nil {
		errOutput := models.ErrorOutput{
			ErrorCode: ErrCodeCouldNotCreateCompany,
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
		Msg("create company executed successfully")
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

	// Name field only has 15 chars, it is not long enough to contain XSS content
	if updateCompanyInput.Description != nil {
		err = xss.CheckForXSS(*updateCompanyInput.Description)
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
				Msg("error while checking for XSS content")
			c.JSON(http.StatusBadRequest, errOutput)
			return
		}
	}

	companyOutput, err := handler.service.PatchCompany(ctx, companyId, updateCompanyInput)
	if err != nil {
		errOutput := models.ErrorOutput{
			ErrorCode: ErrCodePatchCompany,
		}
		err = errors.Join(ErrPatchCompany, err)
		errorCode := http.StatusInternalServerError
		if errors.Is(err, mongo.ErrNoDocuments) {
			errorCode = http.StatusNotFound
		}
		log.Error().
			Err(err).
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Int(consts.LogKeyErrorCode, errOutput.ErrorCode).
			Int(consts.LogKeyStatusCode, errorCode).
			Str(consts.LogKeyCompanyId, companyId.String()).
			Msg("error while trying to PATCH company")
		c.JSON(errorCode, errOutput)
		return
	}

	log.Info().
		Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
		Int(consts.LogKeyStatusCode, http.StatusAccepted).
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

	err = handler.service.DeleteCompany(ctx, companyId)
	if err != nil {
		errOutput := models.ErrorOutput{
			ErrorCode: ErrCodeDeleteCompany,
		}
		err = errors.Join(ErrDeleteCompany, err)
		statusCode := http.StatusInternalServerError
		if errors.Is(err, repo.ErrDocumentNotFound) {
			statusCode = http.StatusNotFound
		}
		log.Error().
			Err(err).
			Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
			Int(consts.LogKeyErrorCode, errOutput.ErrorCode).
			Int(consts.LogKeyStatusCode, statusCode).
			Str(consts.LogKeyCompanyId, companyId.String()).
			Msg("error while trying to delete company")
		c.JSON(statusCode, errOutput)
		return
	}

	log.Info().
		Str(consts.LogKeyTimeUTC, time.Now().UTC().String()).
		Int(consts.LogKeyStatusCode, http.StatusNoContent).
		Str(consts.LogKeyCompanyId, companyId.String()).
		Msg("delete company executed successfully")
	c.JSON(http.StatusNoContent, nil)
}
