package request

type CreateOrderRequest struct {
	ProductID uint   `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
	BuyerID   string `json:"buyer_id" validate:"required"`
}

type CreateSettlementJobRequest struct {
	From string `json:"from" validate:"required"`
	To   string `json:"to" validate:"required"`
}
