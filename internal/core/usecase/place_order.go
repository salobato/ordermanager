package usecase

import (
	"context"
	"errors"

	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/salobato/ordermanager/internal/core/publisher"
	"github.com/salobato/ordermanager/internal/core/repository"
)

type PlaceOrderInput struct {
	CustomerID string
	Total      float64
}

func (i PlaceOrderInput) Validate() error {
	if i.CustomerID == "" {
		return errors.New("O ID do cliente não pode ser vazio")
	}
	if i.Total <= 0 {
		return errors.New("O total deve ser maior que zero")
	}
	return nil
}

func PlaceOrder(r repository.OrderRepository, p publisher.EventPublisher, input PlaceOrderInput) (*entity.Order, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	seq, err := r.NextSequence()
	if err != nil {
		return nil, err
	}

	entry, err := entity.NewOrder(input.CustomerID, input.Total, seq)
	if err != nil {
		return nil, err
	}

	order, err := r.Save(entry)
	if err != nil {
		return nil, err
	}

	err = p.PublishOrderStatusChanged(context.Background(), entity.OrderEvent{
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber.String(),
		Status:      order.Status,
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}
