package usecase_test

import (
	"errors"
	"testing"

	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/salobato/ordermanager/internal/core/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateOrderStatus_WithTestify(t *testing.T) {
	tests := []struct {
		name       string
		fromStatus entity.OrderStatus
		toStatus   entity.OrderStatus
	}{
		{"criado -> em_processamento", entity.OrderCreated, entity.OrderProcessing},
		{"em_processamento -> enviado", entity.OrderProcessing, entity.OrderShipped},
		{"enviado -> entregue", entity.OrderShipped, entity.OrderDelivered},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockOrderRepository)

			order := &entity.Order{
				ID:     "order_id",
				Status: tt.fromStatus,
			}

			repo.On("FindByID", "order_id").Return(order, nil)

			repo.On("UpdateStatus", "order_id", string(tt.toStatus)).
				Return(nil)

			input := usecase.UpdateOrderStatusInput{
				OrderID: "order_id",
				Status:  string(tt.toStatus),
			}

			updatedOrder, err := usecase.UpdateOrderStatus(repo, input)

			assert.NoError(t, err)
			assert.NotNil(t, updatedOrder)
			assert.Equal(t, tt.toStatus, updatedOrder.Status)

			repo.AssertExpectations(t)
		})
	}
}

func TestUpdateOrderStatus_InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input usecase.UpdateOrderStatusInput
	}{
		{
			name:  "ID do pedido vazio",
			input: usecase.UpdateOrderStatusInput{OrderID: "", Status: "enviado"},
		},
		{
			name:  "Status inválido",
			input: usecase.UpdateOrderStatusInput{OrderID: "order_id", Status: "inválido"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockOrderRepository)

			order, err := usecase.UpdateOrderStatus(repo, tt.input)

			assert.Error(t, err)
			assert.Nil(t, order)

			repo.AssertNotCalled(t, "FindByID", mock.Anything)
			repo.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything)
		})
	}
}

func TestUpdateOrderStatus_OrderNotFound(t *testing.T) {
	repo := new(MockOrderRepository)

	repo.On("FindByID", "order_id").
		Return(nil, errors.New("Pedido não encontrado"))

	input := usecase.UpdateOrderStatusInput{
		OrderID: "order_id",
		Status:  string(entity.OrderProcessing),
	}

	order, err := usecase.UpdateOrderStatus(repo, input)

	assert.Error(t, err)
	assert.Nil(t, order)

	repo.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything)
}

func TestUpdateOrderStatus_InvalidTransition(t *testing.T) {
	repo := new(MockOrderRepository)

	order := &entity.Order{
		ID:     "order_id",
		Status: entity.OrderCreated,
	}

	repo.On("FindByID", "order_id").Return(order, nil)

	input := usecase.UpdateOrderStatusInput{
		OrderID: "order_id",
		Status:  string(entity.OrderDelivered),
	}

	updatedOrder, err := usecase.UpdateOrderStatus(repo, input)

	assert.Error(t, err)
	assert.Nil(t, updatedOrder)

	repo.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything)
}

func TestUpdateOrderStatus_UpdateError(t *testing.T) {
	repo := new(MockOrderRepository)

	order := &entity.Order{
		ID:     "order_id",
		Status: entity.OrderProcessing,
	}

	repo.On("FindByID", "order_id").Return(order, nil)

	repo.On("UpdateStatus", "order_id", string(entity.OrderShipped)).
		Return(errors.New("Erro no banco de dados"))

	input := usecase.UpdateOrderStatusInput{
		OrderID: "order_id",
		Status:  string(entity.OrderShipped),
	}

	updatedOrder, err := usecase.UpdateOrderStatus(repo, input)

	assert.Error(t, err)
	assert.Nil(t, updatedOrder)
}
