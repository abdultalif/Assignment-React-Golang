package repository

import (
	"backend-service/internal/core/domain/entity"
	"backend-service/internal/core/domain/model"
	"context"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SettlementRepositoryInterface interface {
	UpsertBatch(ctx context.Context, settlements []entity.SettlementEntity) error
}

type SettlementRepository struct {
	db *gorm.DB
}

// UpsertBatch implements SettlementRepositoryInterface.
func (s *SettlementRepository) UpsertBatch(ctx context.Context, settlements []entity.SettlementEntity) error {

	if len(settlements) == 0 {
		log.Info().Msg("[SettlementRepository] UpsertBatch: no settlements to upsert")
		return nil
	}

	models := make([]model.SettlementModel, len(settlements))
	for i, settlement := range settlements {
		models[i] = model.SettlementModel{
			MerchantID:  settlement.MerchantID,
			Date:        settlement.Date,
			GrossCents:  settlement.GrossCents,
			FeeCents:    settlement.FeeCents,
			NetCents:    settlement.NetCents,
			TxnCount:    settlement.TxnCount,
			GeneratedAt: settlement.GeneratedAt,
			UniqueRunID: settlement.UniqueRunID,
		}
	}

	err := s.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "merchant_id"}, {Name: "date"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"gross_cents", "fee_cents", "net_cents", "txn_count",
				"generated_at", "unique_run_id", "updated_at",
			}),
		}).
		CreateInBatches(&models, 1000).
		Error

	if err != nil {
		log.Error().Err(err).Msg("[SettlementRepository] UpsertBatch: failed to upsert settlements")
		return err
	}

	log.Info().Int("count", len(settlements)).Msg("[SettlementRepository] UpsertBatch: successfully upserted settlements")
	return nil

}

func NewSettlementRepository(db *gorm.DB) SettlementRepositoryInterface {
	return &SettlementRepository{
		db: db,
	}
}
