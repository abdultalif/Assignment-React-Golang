package repository

import (
	"backend-service/internal/core/domain/entity"
	"backend-service/internal/core/domain/model"
	"context"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ProductRepositoryInterface interface {
	GetByID(ctx context.Context, id uint) (*entity.ProductEntity, error)
	UpdateStock(ctx context.Context, id uint, quantity int) error
}

type ProductRepository struct {
	db *gorm.DB
}

// UpdateStock implements ProductRepositoryInterface.
func (p *ProductRepository) UpdateStock(ctx context.Context, id uint, quantity int) error {

	result := p.db.WithContext(ctx).
		Model(&model.ProductModel{}).
		Where("id = ? AND stock >= ?", id, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("failed to update stock")
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Error().Msg("out of stock")
		return gorm.ErrRecordNotFound
	}

	return nil

}

// GetByID implements ProductRepositoryInterface.
func (p *ProductRepository) GetByID(ctx context.Context, productID uint) (*entity.ProductEntity, error) {

	var productModel model.ProductModel
	if err := p.db.WithContext(ctx).First(&productModel, productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Error().Err(err).Msg("product not found")
			return nil, gorm.ErrRecordNotFound
		}
		log.Error().Err(err).Msg("failed to get product")
		return nil, err
	}

	return &entity.ProductEntity{
		ID:         productModel.ID,
		Name:       productModel.Name,
		PriceCents: productModel.PriceCents,
		Stock:      productModel.Stock,
		CreatedAt:  productModel.CreatedAt,
		UpdatedAt:  productModel.UpdatedAt,
	}, nil
}

func NewProductRepository(db *gorm.DB) ProductRepositoryInterface {
	return &ProductRepository{
		db: db,
	}
}
