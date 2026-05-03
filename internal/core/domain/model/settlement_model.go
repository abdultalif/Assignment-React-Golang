package model

import (
	"time"

	"github.com/google/uuid"
)

type SettlementModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MerchantID  string    `gorm:"not null;index:idx_merchant_date"`
	Date        time.Time `gorm:"type:date;not null;index:idx_merchant_date"`
	GrossCents  int64     `gorm:"not null;default:0"`
	FeeCents    int64     `gorm:"not null;default:0"`
	NetCents    int64     `gorm:"not null;default:0"`
	TxnCount    int       `gorm:"not null;default:0"`
	GeneratedAt time.Time
	UniqueRunID string `gorm:"not null;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (SettlementModel) TableName() string {
	return "settlements"
}
