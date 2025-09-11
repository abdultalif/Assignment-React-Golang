package repository

import (
	"backend-service/internal/core/domain/entity"
	errs "backend-service/internal/core/domain/error"
	"backend-service/internal/core/domain/model"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type JobRepositoryInterface interface {
	Create(ctx context.Context, job *entity.JobEntity) (uuid.UUID, error)
	GetByID(ctx context.Context, jobID uuid.UUID) (*entity.JobEntity, error)
}

type JobRepository struct {
	db *gorm.DB
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
