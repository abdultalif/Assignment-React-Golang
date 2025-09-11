package model

import (
	"time"

	"github.com/google/uuid"
)

type JobModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Type         string    `gorm:"not null;index"`
	Status       string    `gorm:"not null;default:QUEUED;index"`
	Progress     int       `gorm:"default:0"`
	Processed    int64     `gorm:"default:0"`
	Total        int64     `gorm:"default:0"`
	Params       string    `gorm:"type:jsonb"`
	ResultPath   *string
	ErrorMessage *string
	UniqueRunID  *string
	Cancelled    bool `gorm:"default:false"`
	StartedAt    *time.Time
	CompletedAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (JobModel) TableName() string {
	return "jobs"
}
