package entity

import (
	"time"

	"github.com/google/uuid"
)

type JobEntity struct {
	ID           uuid.UUID
	Type         string
	Status       string
	Progress     int
	Processed    int64
	Total        int64
	Params       string
	ResultPath   *string
	ErrorMessage *string
	UniqueRunID  *string
	Cancelled    bool
	StartedAt    *time.Time
	CompletedAt  *time.Time
}

type SettlementJobParams struct {
	From string
	To   string
}

type SettlementJob struct {
	ID        uuid.UUID
	From      time.Time
	To        time.Time
	RunID     string
	BatchSize int
	Cancelled chan bool
}
