package service

import (
	"backend-service/internal/adapter/repository"
	"backend-service/internal/core/domain/entity"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type WorkerPool struct {
	jobQueue        chan entity.SettlementJob
	workerCount     int
	transactionRepo repository.TransactionRepositoryInterface
	settlementRepo  repository.SettlementRepositoryInterface
	jobRepo         repository.JobRepositoryInterface
}

func NewWorkerPool(
	workerCount int,
	transactionRepo repository.TransactionRepositoryInterface,
	settlementRepo repository.SettlementRepositoryInterface,
	jobRepo repository.JobRepositoryInterface,
) *WorkerPool {
	return &WorkerPool{
		jobQueue:        make(chan entity.SettlementJob, 100),
		workerCount:     workerCount,
		transactionRepo: transactionRepo,
		settlementRepo:  settlementRepo,
		jobRepo:         jobRepo,
	}
}

func (w *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < w.workerCount; i++ {
		go w.worker(ctx, i)
	}
	log.Info().Int("workers", w.workerCount).Msg("Settlement worker pool started")
}

func (w *WorkerPool) AddJob(job entity.SettlementJob) {
	w.jobQueue <- job
}

func (w *WorkerPool) worker(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Int("worker", workerID).Msg("Worker stopped")
			return
		case job := <-w.jobQueue:
			log.Info().Str("job_id", job.ID.String()).Int("worker", workerID).Msg("Processing settlement job")
			w.processSettlementJob(ctx, job)
		}
	}
}

func (w *WorkerPool) processSettlementJob(ctx context.Context, job entity.SettlementJob) {

	err := w.jobRepo.UpdateStatus(ctx, job.ID, "RUNNING", nil)
	if err != nil {
		log.Error().Err(err).Str("job_id", job.ID.String()).Msg("Failed to update job status to RUNNING")
		return
	}

	now := time.Now()
	err = w.jobRepo.UpdateStartedAt(ctx, job.ID, &now)
	if err != nil {
		log.Error().Err(err).Str("job_id", job.ID.String()).Msg("Failed to update started_at")
	}

	const batchSize = 10000
	var processed int64 = 0

	total, err := w.transactionRepo.Count(ctx, job.From, job.To)
	if err != nil {
		log.Error().Err(err).Str("job_id", job.ID.String()).Msg("Failed to count total transactions")
		w.markJobAsFailed(ctx, job.ID, err.Error())
		return
	}

	settlementsMap := make(map[string]*entity.SettlementEntity)

	for offset := int64(0); offset < total; offset += batchSize {
		select {
		case <-job.Cancelled:
			log.Info().Str("job_id", job.ID.String()).Msg("Job cancelled")
			w.markJobAsCancelled(ctx, job.ID)
			return
		case <-ctx.Done():
			return
		default:
		}

		transactions, err := w.transactionRepo.GetBatch(ctx, job.From, job.To, offset, batchSize)
		if err != nil {
			log.Error().Err(err).Str("job_id", job.ID.String()).Msg("Failed to get transaction batch")
			w.markJobAsFailed(ctx, job.ID, err.Error())
			return
		}

		for _, txn := range transactions {
			dateKey := txn.PaidAt.Format("2006-01-02")
			key := fmt.Sprintf("%s_%s", txn.MerchantID, dateKey)

			if settlement, exists := settlementsMap[key]; exists {
				settlement.GrossCents += int64(txn.AmountCents)
				settlement.FeeCents += int64(txn.FeeCents)
				settlement.NetCents += int64(txn.AmountCents) - int64(txn.FeeCents)
				settlement.TxnCount++
			} else {
				date, _ := time.Parse("2006-01-02", dateKey)
				settlementsMap[key] = &entity.SettlementEntity{
					MerchantID:  txn.MerchantID,
					Date:        date,
					GrossCents:  int64(txn.AmountCents),
					FeeCents:    int64(txn.FeeCents),
					NetCents:    int64(txn.AmountCents) - int64(txn.FeeCents),
					TxnCount:    1,
					GeneratedAt: time.Now(),
					UniqueRunID: job.RunID,
				}
			}
		}

		processed += int64(len(transactions))

		progress := int((processed * 100) / total)
		err = w.jobRepo.UpdateProgress(ctx, job.ID, progress, processed)
		if err != nil {
			log.Error().Err(err).Str("job_id", job.ID.String()).Msg("Failed to update progress")
		}

		log.Info().
			Str("job_id", job.ID.String()).
			Int64("processed", processed).
			Int64("total", total).
			Int("progress", progress).
			Msg("Batch processed")
	}

	settlements := make([]entity.SettlementEntity, 0, len(settlementsMap))
	for _, settlement := range settlementsMap {
		settlements = append(settlements, *settlement)
	}

	err = w.settlementRepo.UpsertBatch(ctx, settlements)
	if err != nil {
		log.Error().Err(err).Str("job_id", job.ID.String()).Msg("Failed to upsert settlements")
		w.markJobAsFailed(ctx, job.ID, err.Error())
		return
	}

	csvPath, err := w.generateCSV(job.ID, settlements)
	if err != nil {
		log.Error().Err(err).Str("job_id", job.ID.String()).Msg("Failed to generate CSV")
		w.markJobAsFailed(ctx, job.ID, err.Error())
		return
	}

	completedAt := time.Now()
	err = w.jobRepo.Complete(ctx, job.ID, csvPath, &completedAt)
	if err != nil {
		log.Error().Err(err).Str("job_id", job.ID.String()).Msg("Failed to mark job as completed")
		return
	}

	log.Info().
		Str("job_id", job.ID.String()).
		Str("csv_path", csvPath).
		Int("settlements_count", len(settlements)).
		Msg("Settlement job completed successfully")
}

func (w *WorkerPool) generateCSV(jobID uuid.UUID, settlements []entity.SettlementEntity) (string, error) {
	dir := "tmp/settlements"
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	filename := fmt.Sprintf("%s.csv", jobID.String())
	filepath := filepath.Join(dir, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"merchant_id", "date", "gross", "fee", "net", "txn_count"}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	for _, settlement := range settlements {
		record := []string{
			settlement.MerchantID,
			settlement.Date.Format("2006-01-02"),
			strconv.FormatInt(settlement.GrossCents, 10),
			strconv.FormatInt(settlement.FeeCents, 10),
			strconv.FormatInt(settlement.NetCents, 10),
			strconv.Itoa(settlement.TxnCount),
		}
		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return filepath, nil
}

func (w *WorkerPool) markJobAsFailed(ctx context.Context, jobID uuid.UUID, errorMsg string) {
	completedAt := time.Now()
	err := w.jobRepo.UpdateStatus(ctx, jobID, "FAILED", &errorMsg)
	if err != nil {
		log.Error().Err(err).Str("job_id", jobID.String()).Msg("Failed to mark job as failed")
	}

	err = w.jobRepo.UpdateCompletedAt(ctx, jobID, &completedAt)
	if err != nil {
		log.Error().Err(err).Str("job_id", jobID.String()).Msg("Failed to update completed_at for failed job")
	}
}

func (w *WorkerPool) markJobAsCancelled(ctx context.Context, jobID uuid.UUID) {
	completedAt := time.Now()
	err := w.jobRepo.UpdateStatus(ctx, jobID, "CANCELLED", nil)
	if err != nil {
		log.Error().Err(err).Str("job_id", jobID.String()).Msg("Failed to mark job as cancelled")
	}

	err = w.jobRepo.UpdateCancelledFlag(ctx, jobID, true)
	if err != nil {
		log.Error().Err(err).Str("job_id", jobID.String()).Msg("Failed to update cancelled flag")
	}

	err = w.jobRepo.UpdateCompletedAt(ctx, jobID, &completedAt)
	if err != nil {
		log.Error().Err(err).Str("job_id", jobID.String()).Msg("Failed to update completed_at for cancelled job")
	}
}
