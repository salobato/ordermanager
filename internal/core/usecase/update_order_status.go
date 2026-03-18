package usecase

import (
	"fmt"

	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/salobato/ordermanager/internal/core/repository"
)

type UpdateOrderStatusInput struct {
	OrderID string
	Status  string
}

func UpdateOrderStatus(repo repository.OrderRepository, input UpdateOrderStatusInput) (*entity.Order, error) {
	if input.OrderID == "" {
		return nil, fmt.Errorf("O ID do pedido não pode ser vazio")
	}

	newStatus := entity.OrderStatus(input.Status)
	if !newStatus.IsValid() {
		return nil, fmt.Errorf("Status do pedido inválido: %s", input.Status)
	}

	order, err := repo.FindByID(input.OrderID)
	if err != nil {
		return nil, fmt.Errorf("Pedido não encontrado: %w", err)
	}

	if err := order.ChangeStatus(newStatus); err != nil {
		return nil, err
	}

	if err := repo.UpdateStatus(order.ID, string(order.Status)); err != nil {
		return nil, err
	}

	return order, nil
}
