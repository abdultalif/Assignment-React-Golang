package request

type CreateOrderRequest struct {
	ProductID uint   `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
	BuyerID   string `json:"buyer_id" validate:"required"`
}
