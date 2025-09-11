package repository

import (
	"backend-service/internal/core/domain/entity"
	errs "backend-service/internal/core/domain/error"
	"backend-service/internal/core/domain/model"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	Create(ctx context.Context, order entity.OrderEntity) (uuid.UUID, error)
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (*entity.OrderEntity, error)
}
type OrderRepository struct {
	db *gorm.DB
}

// GetOrderByID implements OrderRepositoryInterface.
func (o *OrderRepository) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*entity.OrderEntity, error) {

	orderModel := model.OrderModel{}
	if err := o.db.WithContext(ctx).Preload("Product").First(&orderModel, "id = ?", orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msg("order not found")
			return nil, errs.ErrOrderNotFound
		}
		log.Error().Err(err).Msg("failed to get order")
		return nil, err
	}

	return &entity.OrderEntity{
		ID:         orderID,
		ProductID:  orderModel.ProductID,
		BuyerID:    orderModel.BuyerID,
		Quantity:   orderModel.Quantity,
		TotalCents: orderModel.TotalCents,
		Status:     orderModel.Status,
		Product: &entity.ProductEntity{
			ID:         orderModel.Product.ID,
			Name:       orderModel.Product.Name,
			PriceCents: orderModel.Product.PriceCents,
		},
	}, nil

}

// Create implements OrderRepositoryInterface.
func (o *OrderRepository) Create(ctx context.Context, order entity.OrderEntity) (uuid.UUID, error) {
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
		return uuid.Nil, err
	}

	return modelOrder.ID, nil

}

func NewOrderRepository(db *gorm.DB) OrderRepositoryInterface {
	return &OrderRepository{
		db: db,
	}
}
