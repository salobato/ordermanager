package entity

import (
	"errors"
	"time"
)

type Order struct {
	ID          string
	OrderNumber OrderNumber
	Total       float64
	Status      OrderStatus
	CustomerID  string
	PlacedAt    time.Time
	UpdatedAt   time.Time
}

func (o Order) Validate() error {
	if o.CustomerID == "" {
		return errors.New("O ID do cliente não pode ser vazio")
	}
	if o.Total <= 0 {
		return errors.New("O total deve ser maior que zero")
	}
	return nil
}

func NewOrder(customerID string, total float64, sequence int64) (*Order, error) {
	now := time.Now()
	order := &Order{
		CustomerID:  customerID,
		OrderNumber: NewOrderNumber(sequence),
		Total:       total,
		Status:      OrderCreated,
		PlacedAt:    now,
		UpdatedAt:   now,
	}

	if err := order.Validate(); err != nil {
		return nil, err
	}

	return order, nil
}
