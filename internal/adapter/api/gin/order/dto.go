package order

import (
	"time"

	"github.com/salobato/ordermanager/internal/core/entity"
)

type PlaceOrderRequest struct {
	CustomerID string  `json:"customer_id"`
	Total      float64 `json:"total"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

type OrderResponse struct {
	ID          string  `json:"id"`
	OrderNumber string  `json:"order_number"`
	CustomerID  string  `json:"customer_id"`
	Total       float64 `json:"total"`
	Status      string  `json:"status"`
	PlacedAt    string  `json:"placed_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func toOrderResponse(order *entity.Order) OrderResponse {
	return OrderResponse{
		ID:          order.ID,
		OrderNumber: order.OrderNumber.String(),
		CustomerID:  order.CustomerID,
		Total:       order.Total,
		Status:      string(order.Status),
		PlacedAt:    order.PlacedAt.Format(time.RFC3339),
		UpdatedAt:   order.UpdatedAt.Format(time.RFC3339),
	}
}
