package repository

import (
	"backend-service/internal/core/domain/entity"
	"backend-service/internal/core/domain/model"
	"context"

	// "github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	Create(ctx context.Context, order entity.OrderEntity) error
	// GetByID(ctx context.Context, orderID uuid.UUID) (*entity.OrderEntity, error)
}
type OrderRepository struct {
	db *gorm.DB
}

// Create implements OrderRepositoryInterface.
func (o *OrderRepository) Create(ctx context.Context, order entity.OrderEntity) error {
	modelOrder := model.OrderModel{
		ProductID:  order.ProductID,
		BuyerID:    order.BuyerID,
		Quantity:   order.Quantity,
		TotalCents: order.TotalCents,
		Status:     order.Status,
	}

	if err := o.db.Create(&modelOrder).Error; err != nil {
		log.Error().Err(err).
			Str("buyer_id", order.BuyerID).
			Msg("failed to create order")
		return err
	}

	return nil

}

func NewOrderRepository(db *gorm.DB) OrderRepositoryInterface {
	return &OrderRepository{
		db: db,
	}
}
