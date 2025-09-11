package service

import (
	"backend-service/internal/adapter/repository"
	"backend-service/internal/core/domain/entity"
	errs "backend-service/internal/core/domain/error"
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type JobServiceInterface interface {
	CreateSettlementJob(ctx context.Context, from string, to string) (*entity.JobEntity, error)
}

type JobService struct {
	jobRepo         repository.JobRepositoryInterface
	transactionRepo repository.TransactionRepositoryInterface
}

// CreateSettlementJob implements JobServiceInterface.
func (j *JobService) CreateSettlementJob(ctx context.Context, from string, to string) (*entity.JobEntity, error) {

	fromTime, err := time.Parse("2006-01-02", from)
	if err != nil {
		log.Error().Err(err).Msg("[JobService-1] CreateSettlementJob: failed to parse from date")
		return nil, errs.ErrInvalidDateRange
	}

	toTime, err := time.Parse("2006-01-02", to)
	if err != nil {
		log.Error().Err(err).Msg("[JobService-3] CreateSettlementJob: failed to parse from date")
		return nil, errs.ErrInvalidDateRange
	}

	if fromTime.After(toTime) {
		log.Error().Err(err).Msg("[JobService-2] CreateSettlementJob: from date is after to date")
		return nil, errs.ErrInvalidDateRange
	}

	total, err := j.transactionRepo.Count(ctx, fromTime, toTime)
	if err != nil {
		log.Error().Err(err).Msg("[JobService-4] CreateSettlementJob: failed to count transactions")
		return nil, err
	}

	params := entity.SettlementJobParams{
		From: from,
		To:   to,
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		log.Error().Err(err).Msg("[JobService-5] CreateSettlementJob: failed to marshal params")
		return nil, err
	}

	uniqueRunID := uuid.New().String()
	job := &entity.JobEntity{
		Type:        "SETTLEMENT",
		Status:      "QUEUED",
		Total:       total,
		Params:      string(paramsJSON),
		UniqueRunID: &uniqueRunID,
	}

	jobID, err := j.jobRepo.Create(ctx, job)
	if err != nil {
		log.Error().Err(err).Msg("[JobService-6] CreateSettlementJob: failed to create job")
		return nil, err
	}

	job.ID = jobID

	return job, nil

}

func NewJobService(jobRepo repository.JobRepositoryInterface, transactionRepo repository.TransactionRepositoryInterface) JobServiceInterface {
	return &JobService{
		jobRepo:         jobRepo,
		transactionRepo: transactionRepo,
	}
}
