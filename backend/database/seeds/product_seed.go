package seeds

import (
	"backend-service/internal/core/domain/model"
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type SeederInterface interface {
	SeedProducts(ctx context.Context) error
	SeedTransactions(ctx context.Context, count int) error
}

type Seeder struct {
	db *gorm.DB
}

// SeedProductsx implements SeederInterface.
func (s *Seeder) SeedProducts(ctx context.Context) error {
	log.Println("Seeding products...")

	var count int64
	if err := s.db.Model(&model.ProductModel{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("Products already exist, skipping seed")
		return nil
	}

	product := &model.ProductModel{
		Name:       "Limited Edition Product",
		PriceCents: 10000,
		Stock:      100,
	}

	if err := s.db.Create(product).Error; err != nil {
		return fmt.Errorf("failed to seed product: %w", err)
	}

	log.Println("Products seeded successfully")
	return nil
}

func (s *Seeder) SeedTransactions(ctx context.Context, count int) error {
	log.Printf("Seeding %d transactions...", count)

	var existingCount int64
	if err := s.db.Model(&model.TransactionModel{}).Count(&existingCount).Error; err != nil {
		return err
	}

	if existingCount >= int64(count) {
		log.Printf("Transactions already exist (%d), skipping seed", existingCount)
		return nil
	}

	merchants := []string{
		"merchant_001", "merchant_002", "merchant_003", "merchant_004", "merchant_005",
		"merchant_006", "merchant_007", "merchant_008", "merchant_009", "merchant_010",
		"merchant_011", "merchant_012", "merchant_013", "merchant_014", "merchant_015",
		"merchant_016", "merchant_017", "merchant_018", "merchant_019", "merchant_020",
	}

	rand.Seed(time.Now().UnixNano())

	batchSize := 10000
	batches := count / batchSize
	if count%batchSize != 0 {
		batches++
	}

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)
	dateRange := endDate.Sub(startDate)

	for batch := 0; batch < batches; batch++ {
		currentBatchSize := batchSize
		if batch == batches-1 && count%batchSize != 0 {
			currentBatchSize = count % batchSize
		}

		transactions := make([]*model.TransactionModel, currentBatchSize)

		for i := 0; i < currentBatchSize; i++ {

			randomDuration := time.Duration(rand.Int63n(int64(dateRange)))
			paidAt := startDate.Add(randomDuration)

			merchantID := merchants[rand.Intn(len(merchants))]

			amountCents := rand.Intn(49900) + 100

			feeCents := int(float64(amountCents) * 0.03)

			statusRand := rand.Intn(100)
			var status string
			switch {
			case statusRand < 90:
				status = "PAID"
			case statusRand < 95:
				status = "FAILED"
			default:
				status = "PENDING"
			}

			transactions[i] = &model.TransactionModel{
				MerchantID:  merchantID,
				AmountCents: amountCents,
				FeeCents:    feeCents,
				Status:      status,
				PaidAt:      paidAt,
			}
		}

		if err := s.db.CreateInBatches(transactions, batchSize).Error; err != nil {
			return fmt.Errorf("failed to seed transaction batch %d: %w", batch, err)
		}

		log.Printf("Seeded batch %d/%d (%d transactions)", batch+1, batches, currentBatchSize)
	}

	log.Printf("Successfully seeded %d transactions", count)
	return nil
}

func NewSeeder(db *gorm.DB) SeederInterface {
	return &Seeder{db: db}
}
