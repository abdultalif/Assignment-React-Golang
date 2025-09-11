package repository

import (
	"backend-service/internal/core/domain/entity"
	"backend-service/internal/core/domain/model"
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type JobRepositoryInterface interface {
	Create(ctx context.Context, job *entity.JobEntity) (uuid.UUID, error)
}

type JobRepository struct {
	db *gorm.DB
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
