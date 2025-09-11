package response

import "github.com/google/uuid"

type CreateOrderResponse struct {
	OrderID uuid.UUID `json:"order_id"`
	Status  string    `json:"status"`
}

type GetOrderByIDResponse struct {
	OrderID    uuid.UUID       `json:"order_id"`
	ProductID  uint            `json:"product_id"`
	BuyerID    string          `json:"buyer_id"`
	Quantity   int             `json:"quantity"`
	TotalCents int             `json:"total_cents"`
	Status     string          `json:"status"`
	Product    *ProductDetails `json:"product,omitempty"`
}

type ProductDetails struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	PriceCents int    `json:"price_cents"`
}

type CreateJobResponse struct {
	JobID  uuid.UUID `json:"job_id"`
	Status string    `json:"status"`
}
