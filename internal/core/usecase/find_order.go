package usecase

import (
	"errors"
	"fmt"

	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/salobato/ordermanager/internal/core/repository"
)

type FindOrderByIDInput struct {
	OrderID string
}

func FindOrderByID(
	repo repository.OrderRepository,
	input FindOrderByIDInput,
) (*entity.Order, error) {

	if input.OrderID == "" {
		return nil, errors.New("O ID do pedido não pode ser vazio")
	}

	order, err := repo.FindByID(input.OrderID)
	if err != nil {
		return nil, fmt.Errorf("Erro ao buscar pedido: %w", err)
	}

	return order, nil
}
