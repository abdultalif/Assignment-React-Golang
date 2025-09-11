package handler

import (
	"backend-service/internal/adapter/handler/request"
	"backend-service/internal/adapter/handler/response"
	errs "backend-service/internal/core/domain/error"
	"backend-service/internal/core/service"
	v "backend-service/pkg/validator"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type JobHandlerInterface interface {
	CreateSettlementJob(c *gin.Context)
}

type JobHandler struct {
	jobService service.JobServiceInterface
	validator  *v.Validator
}

// CreateSettlementJob implements JobHandlerInterface.
func (j *JobHandler) CreateSettlementJob(c *gin.Context) {

	var (
		ctx = c.Request.Context()
		req = request.CreateSettlementJobRequest{}
		res = response.CreateJobResponse{}
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("[OrderHandler-1] CreateSettlementJob")
		c.JSON(http.StatusBadRequest, response.ResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := j.validator.Validate(req); err != nil {

		log.Error().Err(err).Msg("[OrderHandler-2] CreateSettlementJob")

		if ve, ok := err.(v.ValidationError); ok {
			c.JSON(http.StatusUnprocessableEntity, response.ResponseError(http.StatusUnprocessableEntity, ve.Errors))
			return
		}

		c.JSON(http.StatusUnprocessableEntity, response.ResponseError(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	job, err := j.jobService.CreateSettlementJob(ctx, req.From, req.To)
	if err != nil {
		log.Error().Err(err).Msg("[OrderHandler-3] CreateSettlementJob")
		if errors.Is(err, errs.ErrInvalidDateRange) {
			c.JSON(http.StatusBadRequest, response.ResponseError(http.StatusBadRequest, err.Error()))
			return
		} else {
			c.JSON(http.StatusInternalServerError, response.ResponseError(http.StatusInternalServerError, err.Error()))
			return
		}
	}

	res.JobID = job.ID
	res.Status = job.Status

	c.JSON(http.StatusAccepted, response.ResponseSuccess(http.StatusAccepted, "success", res))

}

func NewJobHandler(jobService service.JobServiceInterface, validator *v.Validator) JobHandlerInterface {
	return &JobHandler{
		jobService: jobService,
		validator:  validator,
	}
}
