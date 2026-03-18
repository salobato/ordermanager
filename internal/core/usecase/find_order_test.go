package usecase_test

import (
	"errors"
	"testing"

	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/salobato/ordermanager/internal/core/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindOrderByID_WithTestify(t *testing.T) {
	repo := new(MockOrderRepository)

	expectedOrder := &entity.Order{
		ID:         "order_id",
		CustomerID: "customer_id",
		Total:      100,
		Status:     entity.OrderCreated,
	}

	repo.On("FindByID", "order_id").Return(expectedOrder, nil)

	input := usecase.FindOrderByIDInput{
		OrderID: "order_id",
	}

	order, err := usecase.FindOrderByID(repo, input)

	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, "order_id", order.ID)

	repo.AssertExpectations(t)
}

func TestFindOrderByID_InvalidInput(t *testing.T) {
	repo := new(MockOrderRepository)

	input := usecase.FindOrderByIDInput{
		OrderID: "",
	}

	order, err := usecase.FindOrderByID(repo, input)

	assert.Error(t, err)
	assert.Nil(t, order)

	repo.AssertNotCalled(t, "FindByID", mock.Anything)
}

func TestFindOrderByID_NotFound(t *testing.T) {
	repo := new(MockOrderRepository)

	repo.On("FindByID", "order_id").
		Return(nil, errors.New("Pedido não encontrado"))

	input := usecase.FindOrderByIDInput{
		OrderID: "order_id",
	}

	order, err := usecase.FindOrderByID(repo, input)

	assert.Error(t, err)
	assert.Nil(t, order)
}

func TestFindOrderByID_RepositoryError(t *testing.T) {
	repo := new(MockOrderRepository)

	repo.On("FindByID", "order_id").
		Return(nil, errors.New("db error"))

	input := usecase.FindOrderByIDInput{
		OrderID: "order_id",
	}

	order, err := usecase.FindOrderByID(repo, input)

	assert.Error(t, err)
	assert.Nil(t, order)
}
