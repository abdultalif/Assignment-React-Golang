package model

import "time"

type ProductModel struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
	PriceCents int    `gorm:"not null;default:0"`
	Stock      int    `gorm:"not null;default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (ProductModel) TableName() string {
	return "products"
}
