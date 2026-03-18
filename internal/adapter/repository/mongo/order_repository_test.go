package mongo_test

import (
	"testing"
	"time"

	"github.com/salobato/ordermanager/internal/adapter/repository/mongo"
	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/stretchr/testify/assert"
)

func TestOrderRepository_Save_Insert(t *testing.T) {
	db := setupTestDB(t)

	counterRepo := mongo.NewCounterRepository(db)
	repo := mongo.NewOrderRepository(db, counterRepo)

	order := &entity.Order{
		OrderNumber: "ORD-2026-000001",
		CustomerID:  "507f1f77bcf86cd799439011",
		Total:       100,
		Status:      entity.OrderCreated,
		PlacedAt:    time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := repo.Save(order)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID)
}

func TestOrderRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)

	counterRepo := mongo.NewCounterRepository(db)
	repo := mongo.NewOrderRepository(db, counterRepo)

	order := &entity.Order{
		OrderNumber: "ORD-2026-000001",
		CustomerID:  "507f1f77bcf86cd799439011",
		Total:       100,
		Status:      entity.OrderCreated,
		PlacedAt:    time.Now(),
		UpdatedAt:   time.Now(),
	}

	saved, _ := repo.Save(order)

	found, err := repo.FindByID(saved.ID)

	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, saved.ID, found.ID)
}

func TestOrderRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)

	counterRepo := mongo.NewCounterRepository(db)
	repo := mongo.NewOrderRepository(db, counterRepo)

	order, err := repo.FindByID("507f1f77bcf86cd799439011")

	assert.Error(t, err)
	assert.Nil(t, order)
}

func TestOrderRepository_UpdateStatus(t *testing.T) {
	db := setupTestDB(t)

	counterRepo := mongo.NewCounterRepository(db)
	repo := mongo.NewOrderRepository(db, counterRepo)

	order := &entity.Order{
		OrderNumber: "ORD-2026-000001",
		CustomerID:  "507f1f77bcf86cd799439011",
		Total:       100,
		Status:      entity.OrderCreated,
		PlacedAt:    time.Now(),
		UpdatedAt:   time.Now(),
	}

	saved, _ := repo.Save(order)

	err := repo.UpdateStatus(saved.ID, string(entity.OrderProcessing))

	assert.NoError(t, err)

	updated, _ := repo.FindByID(saved.ID)

	assert.Equal(t, entity.OrderProcessing, updated.Status)
}
