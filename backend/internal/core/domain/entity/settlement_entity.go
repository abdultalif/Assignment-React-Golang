package entity

import (
	"time"

	"github.com/google/uuid"
)

type SettlementEntity struct {
	ID          uuid.UUID
	MerchantID  string
	Date        time.Time
	GrossCents  int64
	FeeCents    int64
	NetCents    int64
	TxnCount    int
	GeneratedAt time.Time
	UniqueRunID string
}
