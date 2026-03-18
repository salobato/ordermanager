package usecase_test

import (
	"errors"
	"testing"

	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/salobato/ordermanager/internal/core/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPlaceOrder_WithTestify(t *testing.T) {
	repo := new(MockOrderRepository)

	repo.On("NextSequence").Return(int64(1), nil)

	repo.On("Save", mock.MatchedBy(func(order *entity.Order) bool {
		return order.CustomerID == "customer_id" &&
			order.Total == 299.90 &&
			order.Status == entity.OrderCreated &&
			order.OrderNumber.IsValid()
	})).Return(&entity.Order{
		ID:          "mongo_id",
		OrderNumber: entity.NewOrderNumber(1),
		CustomerID:  "customer_id",
		Total:       299.90,
		Status:      entity.OrderCreated,
	}, nil)

	input := usecase.PlaceOrderInput{
		CustomerID: "customer_id",
		Total:      299.90,
	}

	order, err := usecase.PlaceOrder(repo, input)

	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, "mongo_id", order.ID)
	assert.True(t, order.OrderNumber.IsValid())

	repo.AssertExpectations(t)
}

func TestPlaceOrder_InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input usecase.PlaceOrderInput
	}{
		{
			name:  "Customer ID vazio",
			input: usecase.PlaceOrderInput{CustomerID: "", Total: 100},
		},
		{
			name:  "Total negativo",
			input: usecase.PlaceOrderInput{CustomerID: "customer_id", Total: -10},
		},
		{
			name:  "Total igual a zero",
			input: usecase.PlaceOrderInput{CustomerID: "customer_id", Total: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockOrderRepository)

			order, err := usecase.PlaceOrder(repo, tt.input)

			assert.Error(t, err)
			assert.Nil(t, order)
			repo.AssertNotCalled(t, "NextSequence")
			repo.AssertNotCalled(t, "Save", mock.Anything)
		})
	}
}

func TestPlaceOrder_SequenceError(t *testing.T) {
	repo := new(MockOrderRepository)

	repo.On("NextSequence").Return(int64(0), errors.New("sequence error"))

	input := usecase.PlaceOrderInput{
		CustomerID: "customer_id",
		Total:      100,
	}

	order, err := usecase.PlaceOrder(repo, input)

	assert.Error(t, err)
	assert.Nil(t, order)

	repo.AssertNotCalled(t, "Save", mock.Anything)
}

func TestPlaceOrder_RepositoryError(t *testing.T) {
	repo := new(MockOrderRepository)

	repo.On("NextSequence").Return(int64(1), nil)
	repo.On("Save", mock.Anything).
		Return((*entity.Order)(nil), errors.New("db error"))

	input := usecase.PlaceOrderInput{
		CustomerID: "customer_id",
		Total:      100,
	}

	order, err := usecase.PlaceOrder(repo, input)

	assert.Error(t, err)
	assert.Nil(t, order)
}
