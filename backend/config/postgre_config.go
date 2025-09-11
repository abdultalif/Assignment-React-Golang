package config

import (
	seeder "backend-service/database/seeds"
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func (cfg Config) ConnectionPostgres() (*Postgres, error) {
	dbConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName)

	db, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		return nil, err
	}

	seeds := seeder.NewSeeder(db)
	if err := seeds.SeedProducts(context.Background()); err != nil {
		log.Error().Err(err).Msg("Failed to seed products")
		return nil, err
	}

	if err := seeds.SeedTransactions(context.Background(), 1000000); err != nil {
		log.Error().Err(err).Msg("Failed to seed transactions")
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.Postgres.DBMaxOpen)
	sqlDB.SetMaxIdleConns(cfg.Postgres.DBMaxIdle)

	return &Postgres{
		DB: db,
	}, nil
}
