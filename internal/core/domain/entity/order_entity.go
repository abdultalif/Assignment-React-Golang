package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderEntity struct {
	ID         uuid.UUID
	ProductID  uint
	BuyerID    string
	Quantity   int
	TotalCents int
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Product    *ProductEntity
}
