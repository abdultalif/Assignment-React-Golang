package seeds

import (
	"backend-service/internal/core/domain/model"
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type SeederInterface interface {
	SeedProducts(ctx context.Context) error
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

func NewSeeder(db *gorm.DB) SeederInterface {
	return &Seeder{db: db}
}
