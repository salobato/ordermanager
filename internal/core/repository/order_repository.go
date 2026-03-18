package repository

import "github.com/salobato/ordermanager/internal/core/entity"

type OrderRepository interface {
	Save(order *entity.Order) (*entity.Order, error)
	FindByID(id string) (*entity.Order, error)
	UpdateStatus(id string, status string) error
	NextSequence() (int64, error)
}
