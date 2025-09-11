package repository

import (
	"backend-service/internal/core/domain/model"
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type TransactionRepositoryInterface interface {
	Count(ctx context.Context, from, to time.Time) (int64, error)
}
type TransactionRepository struct {
	db *gorm.DB
}

// Count implements TransactionRepositoryInterface.
func (t *TransactionRepository) Count(ctx context.Context, from time.Time, to time.Time) (int64, error) {

	var count int64
	err := t.db.WithContext(ctx).
		Model(&model.TransactionModel{}).
		Where("paid_at >= ? AND paid_at <= ?", from, to).
		Count(&count).Error

	if err != nil {
		log.Error().Err(err).Msg("[TransactionRepository-1] Count: failed to count transactions")
		return 0, err
	}

	return count, err
}

func NewTransactionRepository(db *gorm.DB) TransactionRepositoryInterface {
	return &TransactionRepository{db: db}
}
