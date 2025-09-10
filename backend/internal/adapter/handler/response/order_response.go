package response

import "github.com/google/uuid"

type CreateOrderResponse struct {
	OrderID uuid.UUID `json:"order_id"`
	Status  string    `json:"status"`
}
