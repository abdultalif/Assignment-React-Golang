package handler

import (
	"backend-service/internal/adapter/handler/request"
	"backend-service/internal/adapter/handler/response"
	errs "backend-service/internal/core/domain/error"
	"backend-service/internal/core/service"
	v "backend-service/pkg/validator"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type JobHandlerInterface interface {
	CreateSettlementJob(c *gin.Context)
	GetJob(c *gin.Context)
	CancelJob(c *gin.Context)
}

type JobHandler struct {
	jobService service.JobServiceInterface
	validator  *v.Validator
}

// CancelJob implements JobHandlerInterface.
func (j *JobHandler) CancelJob(c *gin.Context) {
	var (
		ctx = c.Request.Context()
	)

	jobID, err := uuid.Parse(c.Param("jobID"))
	if err != nil {
		log.Error().Err(err).Msg("[JobHandler-1] CancelJob: invalid job ID")
		c.JSON(http.StatusBadRequest, response.ResponseError(http.StatusBadRequest, "invalid job ID"))
		return
	}

	err = j.jobService.CancelJob(ctx, jobID)
	if err != nil {
		log.Error().Err(err).Str("job_id", jobID.String()).Msg("[JobHandler-2] CancelJob: failed to cancel job")

		if errors.Is(err, errs.ErrJobNotFound) {
			c.JSON(http.StatusNotFound, response.ResponseError(http.StatusNotFound, err.Error()))
			return
		} else if errors.Is(err, errs.ErrJobCannotBeCancelled) {
			c.JSON(http.StatusBadRequest, response.ResponseError(http.StatusBadRequest, err.Error()))
			return
		} else {
			c.JSON(http.StatusInternalServerError, response.ResponseError(http.StatusInternalServerError, err.Error()))
			return
		}
	}

	c.JSON(http.StatusOK, response.ResponseSuccess(http.StatusOK, "job cancellation requested", nil))
}

// GetJob implements JobHandlerInterface.
func (j *JobHandler) GetJob(c *gin.Context) {

	var (
		ctx = c.Request.Context()
		res = response.JobStatusResponse{}
	)

	jobID, err := uuid.Parse(c.Param("jobID"))
	if err != nil {
		log.Error().Err(err).Msg("[JobHandler-1] GetJob: invalid job ID")
		c.JSON(http.StatusBadRequest, response.ResponseError(http.StatusBadRequest, "invalid job ID"))
		return
	}

	job, err := j.jobService.GetJob(ctx, jobID)
	if err != nil {
		log.Error().Err(err).Msg("[JobHandler-2] GetJob: failed to get job")
		if errors.Is(err, errs.ErrJobNotFound) {
			c.JSON(http.StatusNotFound, response.ResponseError(http.StatusNotFound, err.Error()))
			return
		} else {
			c.JSON(http.StatusInternalServerError, response.ResponseError(http.StatusInternalServerError, err.Error()))
			return
		}
	}

	res.JobID = job.ID
	res.Status = job.Status
	res.Total = job.Total
	res.Progress = job.Progress
	res.Processed = job.Processed

	if job.Status == "COMPLETED" && job.ResultPath != nil {
		downloadURL := "/downloads/" + strings.TrimSuffix(strings.TrimPrefix(*job.ResultPath, "/tmp/settlements/"), ".csv") + ".csv"
		res.DownloadURL = &downloadURL
	}

	c.JSON(http.StatusOK, response.ResponseSuccess(http.StatusOK, "success", res))
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
