package service

import (
	"backend-service/internal/adapter/repository"
	"backend-service/internal/core/domain/entity"
	errs "backend-service/internal/core/domain/error"
	"context"

	"gorm.io/gorm"
)

type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, order entity.OrderEntity) (*entity.OrderEntity, error)
	// CreateOrder(ctx context.Context, productID uint, quantity int, buyerID string) (*entity.Order, error)
}

type OrderService struct {
	orderRepo   repository.OrderRepositoryInterface
	productRepo repository.ProductRepositoryInterface
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

	err = o.orderRepo.Create(ctx, *newOrder)
	if err != nil {
		return nil, err
	}

	return newOrder, nil

}

func NewOrderService() OrderServiceInterface {
	return &OrderService{}
}
