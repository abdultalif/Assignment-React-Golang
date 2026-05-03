package repository

import (
	"backend-service/internal/core/domain/entity"
	"backend-service/internal/core/domain/model"
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type TransactionRepositoryInterface interface {
	Count(ctx context.Context, from, to time.Time) (int64, error)
	GetBatch(ctx context.Context, from time.Time, to time.Time, offset int64, limit int64) ([]entity.TransactionEntity, error)
}

type TransactionRepository struct {
	db *gorm.DB
}

// GetBatch implements TransactionRepositoryInterface.
func (t *TransactionRepository) GetBatch(ctx context.Context, from time.Time, to time.Time, offset int64, limit int64) ([]entity.TransactionEntity, error) {

	var transactions []model.TransactionModel

	err := t.db.WithContext(ctx).
		Where("paid_at >= ? AND paid_at <= ? AND status = ?", from, to, "PAID").
		Offset(int(offset)).
		Limit(int(limit)).
		Order("paid_at ASC, id ASC").
		Find(&transactions).Error

	if err != nil {
		log.Error().Err(err).Msg("[TransactionRepository] GetBatch: failed to get transaction batch")
		return nil, err
	}

	entities := make([]entity.TransactionEntity, len(transactions))
	for i, txn := range transactions {
		entities[i] = entity.TransactionEntity{
			ID:          txn.ID,
			MerchantID:  txn.MerchantID,
			AmountCents: txn.AmountCents,
			FeeCents:    txn.FeeCents,
			Status:      txn.Status,
			PaidAt:      txn.PaidAt,
		}
	}

	return entities, nil

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
