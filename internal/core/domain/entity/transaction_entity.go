package entity

import (
	"time"

	"github.com/google/uuid"
)

type TransactionEntity struct {
	ID          uuid.UUID
	MerchantID  string
	AmountCents int
	FeeCents    int
	Status      string
	PaidAt      time.Time
}
