package repository

import (
	"backend-service/internal/core/domain/entity"
	errs "backend-service/internal/core/domain/error"
	"backend-service/internal/core/domain/model"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type JobRepositoryInterface interface {
	Create(ctx context.Context, job *entity.JobEntity) (uuid.UUID, error)
	GetByID(ctx context.Context, jobID uuid.UUID) (*entity.JobEntity, error)
	UpdateStatus(ctx context.Context, jobID uuid.UUID, status string, errorMessage *string) error
	UpdateStartedAt(ctx context.Context, jobID uuid.UUID, startedAt *time.Time) error
	UpdateProgress(ctx context.Context, jobID uuid.UUID, progress int, processed int64) error
	Complete(ctx context.Context, jobID uuid.UUID, resultPath string, completedAt *time.Time) error
	UpdateCompletedAt(ctx context.Context, jobID uuid.UUID, completedAt *time.Time) error
}

type JobRepository struct {
	db *gorm.DB
}

// UpdateCompletedAt implements JobRepositoryInterface.
func (j *JobRepository) UpdateCompletedAt(ctx context.Context, jobID uuid.UUID, completedAt *time.Time) error {

	err := j.db.WithContext(ctx).
		Model(&model.JobModel{}).
		Where("id = ?", jobID).
		Update("completed_at", completedAt).Error

	if err != nil {
		log.Error().Err(err).Str("job_id", jobID.String()).Msg("[JobRepository] UpdateCompletedAt: failed to update completed_at")
		return err
	}

	return nil

}

// Complete implements JobRepositoryInterface.
func (j *JobRepository) Complete(ctx context.Context, jobID uuid.UUID, resultPath string, completedAt *time.Time) error {

	err := j.db.WithContext(ctx).
		Model(&model.JobModel{}).
		Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"status":       "COMPLETED",
			"progress":     100,
			"result_path":  resultPath,
			"completed_at": completedAt,
		}).Error

	if err != nil {
		log.Error().Err(err).Str("job_id", jobID.String()).Msg("[JobRepository] Complete: failed to mark job as completed")
		return err
	}

	return nil

}

// UpdateProgress implements JobRepositoryInterface.
func (j *JobRepository) UpdateProgress(ctx context.Context, jobID uuid.UUID, progress int, processed int64) error {
	err := j.db.WithContext(ctx).
		Model(&model.JobModel{}).
		Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"progress":  progress,
			"processed": processed,
		}).Error

	if err != nil {
		log.Error().Err(err).Str("job_id", jobID.String()).Msg("[JobRepository] UpdateProgress: failed to update job progress")
		return err
	}

	return nil
}

// UpdateStartedAt implements JobRepositoryInterface.
func (j *JobRepository) UpdateStartedAt(ctx context.Context, jobID uuid.UUID, startedAt *time.Time) error {

	err := j.db.WithContext(ctx).
		Model(&model.JobModel{}).
		Where("id = ?", jobID).
		Update("started_at", startedAt).Error

	if err != nil {
		log.Error().Err(err).Str("job_id", jobID.String()).Msg("[JobRepository] UpdateStartedAt: failed to update started_at")
		return err
	}

	return nil

}

// UpdateStatus implements JobRepositoryInterface.
func (j *JobRepository) UpdateStatus(ctx context.Context, jobID uuid.UUID, status string, errorMessage *string) error {

	updates := map[string]interface{}{
		"status": status,
	}

	if errorMessage != nil {
		updates["error_message"] = *errorMessage
	}

	err := j.db.WithContext(ctx).
		Model(&model.JobModel{}).
		Where("id = ?", jobID).
		Updates(updates).Error

	if err != nil {
		log.Error().Err(err).Str("job_id", jobID.String()).Msg("[JobRepository] UpdateStatus: failed to update job status")
		return err
	}

	return nil

}

// GetByID implements JobRepositoryInterface.
func (j *JobRepository) GetByID(ctx context.Context, jobID uuid.UUID) (*entity.JobEntity, error) {

	modelJob := model.JobModel{}
	if err := j.db.WithContext(ctx).First(&modelJob, "id = ?", jobID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("[JobRepository-1] GetByID: job not found")
			return nil, errs.ErrJobNotFound
		}
		log.Error().Err(err).Msg("[JobRepository-2] GetByID: failed to get job by ID")
		return nil, err
	}
	return &entity.JobEntity{
		ID:          modelJob.ID,
		Type:        modelJob.Type,
		Status:      modelJob.Status,
		Total:       modelJob.Total,
		Params:      modelJob.Params,
		UniqueRunID: modelJob.UniqueRunID,
		ResultPath:  modelJob.ResultPath,
	}, nil

}

// Create implements JobRepositoryInterface.
func (j *JobRepository) Create(ctx context.Context, job *entity.JobEntity) (uuid.UUID, error) {
	request := model.JobModel{
		Type:        job.Type,
		Status:      job.Status,
		Total:       job.Total,
		Params:      job.Params,
		UniqueRunID: job.UniqueRunID,
	}

	if err := j.db.WithContext(ctx).Create(&request).Error; err != nil {
		log.Error().Err(err).Msg("[JobRepository-1] Create: failed to create job")
		return uuid.Nil, err
	}

	return request.ID, nil
}

func NewJobRepository(db *gorm.DB) JobRepositoryInterface {
	return &JobRepository{db: db}
}
