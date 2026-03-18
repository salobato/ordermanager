package entity

import (
	"fmt"
	"time"
)

type OrderStatus string

const (
	OrderCreated    OrderStatus = "criado"
	OrderProcessing OrderStatus = "em_processamento"
	OrderShipped    OrderStatus = "enviado"
	OrderDelivered  OrderStatus = "entregue"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderCreated, OrderProcessing, OrderShipped, OrderDelivered:
		return true
	default:
		return false
	}
}

func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	if s == OrderDelivered {
		return false
	}

	switch s {
	case OrderCreated:
		return next == OrderProcessing
	case OrderProcessing:
		return next == OrderShipped
	case OrderShipped:
		return next == OrderDelivered
	default:
		return false
	}
}

func (o *Order) ChangeStatus(s OrderStatus) error {
	if !o.Status.CanTransitionTo(s) {
		return fmt.Errorf("O pedido não pode trocar de %s para %s", o.Status, s)
	}

	o.Status = s
	o.UpdatedAt = time.Now()

	return nil
}
