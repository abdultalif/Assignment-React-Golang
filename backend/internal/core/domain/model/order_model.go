package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderModel struct {
	ID         uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID  uint         `gorm:"not null"`
	BuyerID    string       `gorm:"not null;index"`
	Quantity   int          `gorm:"not null"`
	TotalCents int          `gorm:"not null"`
	Status     string       `gorm:"default:PENDING"`
	CreatedAt  time.Time    `gorm:"autoCreateTime"`
	UpdatedAt  time.Time    `gorm:"autoUpdateTime"`
	Product    ProductModel `gorm:"foreignKey:ProductID"`
}

func (OrderModel) TableName() string {
	return "orders"
}
