package usecase_test

import (
	"context"

	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/stretchr/testify/mock"
)

type MockEventPublisher struct {
	mock.Mock
}

func (m *MockEventPublisher) PublishOrderStatusChanged(
	ctx context.Context,
	event entity.OrderEvent,
) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

type MockOrderRepository struct {
	mock.Mock
}

type MockCounterRepository struct {
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

func (m *MockCounterRepository) GetNextSequence(counterName string) (int64, error) {
	args := m.Called(counterName)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCounterRepository) GetCurrentSequence(counterName string) (int64, error) {
	args := m.Called(counterName)
	return args.Get(0).(int64), args.Error(1)
}
