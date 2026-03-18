package usecase_test

import (
	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Save(order *entity.Order) (*entity.Order, error) {
	args := m.Called(order)
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByID(id string) (*entity.Order, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *MockOrderRepository) UpdateStatus(id string, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockOrderRepository) NextSequence() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
