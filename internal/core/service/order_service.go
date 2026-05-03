package service

import (
	"backend-service/internal/adapter/repository"
	"backend-service/internal/core/domain/entity"
	errs "backend-service/internal/core/domain/error"
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, order entity.OrderEntity) (*entity.OrderEntity, error)
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (*entity.OrderEntity, error)
}

type OrderService struct {
	orderRepo   repository.OrderRepositoryInterface
	productRepo repository.ProductRepositoryInterface
}

// GetOrderByID implements OrderServiceInterface.
func (o *OrderService) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*entity.OrderEntity, error) {

	order, err := o.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get order by ID")
	}

	return order, err

}

// CreateOrder implements OrderServiceInterface.
func (o *OrderService) CreateOrder(ctx context.Context, order entity.OrderEntity) (*entity.OrderEntity, error) {

	product, err := o.productRepo.GetByID(ctx, order.ProductID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.ErrProductNotFound
		}
		return nil, err
	}

	// Check stock and update atomically
	err = o.productRepo.UpdateStock(ctx, order.ProductID, order.Quantity)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.ErrOutOfStock
		}
		return nil, err
	}

	newOrder := &entity.OrderEntity{
		ProductID:  order.ProductID,
		BuyerID:    order.BuyerID,
		Quantity:   order.Quantity,
		TotalCents: product.PriceCents * order.Quantity,
		Status:     "COMPLETED",
	}

	orderID, err := o.orderRepo.Create(ctx, *newOrder)
	if err != nil {
		return nil, err
	}

	newOrder.ID = orderID

	return newOrder, nil

}

func NewOrderService(orderRepo repository.OrderRepositoryInterface, productRepo repository.ProductRepositoryInterface) OrderServiceInterface {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}
