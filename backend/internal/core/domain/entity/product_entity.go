package entity

import "time"

type ProductEntity struct {
	ID         uint
	Name       string
	PriceCents int
	Stock      int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
