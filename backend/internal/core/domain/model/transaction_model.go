package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MerchantID  string    `gorm:"not null;index"`
	AmountCents int       `gorm:"not null"`
	FeeCents    int       `gorm:"not null"`
	Status      string    `gorm:"default:PAID"`
	PaidAt      time.Time `gorm:"index"`
	CreatedAt   time.Time
}

func (TransactionModel) TableName() string {
	return "transactions"
}
