package service

import (
	"backend-service/config"
	"backend-service/internal/adapter/repository"
	"backend-service/internal/core/domain/entity"
	errs "backend-service/internal/core/domain/error"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type JobServiceInterface interface {
	CreateSettlementJob(ctx context.Context, from string, to string) (*entity.JobEntity, error)
	GetJob(ctx context.Context, jobID uuid.UUID) (*entity.JobEntity, error)
	StartWorkerPool(ctx context.Context)
	CancelJob(ctx context.Context, jobID uuid.UUID) error
}

type JobService struct {
	jobRepo         repository.JobRepositoryInterface
	transactionRepo repository.TransactionRepositoryInterface
	workerPool      *WorkerPool
	activeJobs      map[uuid.UUID]chan bool
}

// CancelJob implements JobServiceInterface.
func (j *JobService) CancelJob(ctx context.Context, jobID uuid.UUID) error {

	job, err := j.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		log.Error().Err(err).Str("job_id", jobID.String()).Msg("[JobService-1] CancelJob: failed to get job")
		return err
	}

	if job.Status != "QUEUED" && job.Status != "RUNNING" {
		return errs.ErrJobCannotBeCancelled
	}

	if cancelChan, exists := j.activeJobs[jobID]; exists {
		select {
		case cancelChan <- true:
			log.Info().Str("job_id", jobID.String()).Msg("[JobService-2] CancelJob: Cancellation signal sent to worker")
		default:
			log.Warn().Str("job_id", jobID.String()).Msg("[JobService-3] CancelJob: Cancellation signal already sent or worker not listening")
		}
		delete(j.activeJobs, jobID)
	}

	if job.Status == "QUEUED" {
		completedAt := time.Now()
		err = j.jobRepo.UpdateStatus(ctx, jobID, "CANCELLED", nil)
		if err != nil {
			log.Error().Err(err).Str("job_id", jobID.String()).Msg(" [JobService-4] CancelJob: Failed to update job status to CANCELLED")
			return err
		}

		err = j.jobRepo.UpdateCancelledFlag(ctx, jobID, true)
		if err != nil {
			log.Error().Err(err).Str("job_id", jobID.String()).Msg("[JobService-5] CancelJob: Failed to update cancelled flag for cancelled job")
		}

		err = j.jobRepo.UpdateCompletedAt(ctx, jobID, &completedAt)
		if err != nil {
			log.Error().Err(err).Str("job_id", jobID.String()).Msg("[JobService-5] CancelJob: Failed to update completed_at for cancelled job")
		}
	}

	return nil
}

// StartWorkerPool implements JobServiceInterface.
func (j *JobService) StartWorkerPool(ctx context.Context) {
	j.workerPool.Start(ctx)
}

// GetJob implements JobServiceInterface.
func (j *JobService) GetJob(ctx context.Context, jobID uuid.UUID) (*entity.JobEntity, error) {
	return j.jobRepo.GetByID(ctx, jobID)
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
		log.Error().Err(err).Msg("[JobService-3] CreateSettlementJob: failed to parse to date")
		return nil, errs.ErrInvalidDateRange
	}

	if fromTime.After(toTime) {
		log.Error().Msg("[JobService-2] CreateSettlementJob: from date is after to date")
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

	cancelChan := make(chan bool, 1)
	j.activeJobs[jobID] = cancelChan

	settlementJob := entity.SettlementJob{
		ID:        jobID,
		From:      fromTime,
		To:        toTime,
		RunID:     uniqueRunID,
		BatchSize: 100,
		Cancelled: cancelChan,
	}

	go func() {
		j.workerPool.AddJob(settlementJob)

		go func() {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				currentJob, err := j.jobRepo.GetByID(context.Background(), jobID)
				if err != nil {
					break
				}

				if currentJob.Status == "COMPLETED" || currentJob.Status == "FAILED" || currentJob.Status == "CANCELLED" {
					delete(j.activeJobs, jobID)
					close(cancelChan)
					break
				}
			}
		}()
	}()

	log.Info().
		Str("job_id", jobID.String()).
		Str("from", from).
		Str("to", to).
		Int64("total", total).
		Msg("Settlement job created and queued")

	return job, nil
}

func NewJobService(cfg *config.Config, jobRepo repository.JobRepositoryInterface, transactionRepo repository.TransactionRepositoryInterface, settlementRepo repository.SettlementRepositoryInterface) JobServiceInterface {

	workerCount := 4

	if envWorkers := cfg.WORKERS.Count; envWorkers != 0 {
		if parsed, err := strconv.Atoi(strconv.Itoa(envWorkers)); err == nil && parsed > 0 {
			workerCount = parsed
		}
	}

	workerPool := NewWorkerPool(workerCount, transactionRepo, settlementRepo, jobRepo)

	return &JobService{
		jobRepo:         jobRepo,
		transactionRepo: transactionRepo,
		workerPool:      workerPool,
		activeJobs:      make(map[uuid.UUID]chan bool),
	}
}
