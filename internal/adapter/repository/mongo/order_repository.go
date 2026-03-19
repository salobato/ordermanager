package mongo

import (
	"context"
	"time"

	"github.com/salobato/ordermanager/internal/adapter/repository/mongo/models"
	"github.com/salobato/ordermanager/internal/core/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type OrderRepository struct {
	collection  *mongo.Collection
	counterRepo *CounterRepository
}

func NewOrderRepository(db *mongo.Database, counterRepo *CounterRepository) *OrderRepository {
	repo := &OrderRepository{
		collection:  db.Collection("orders"),
		counterRepo: counterRepo,
	}

	repo.ensureIndexes(context.Background())

	return repo
}

func (r *OrderRepository) ensureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{bson.E{Key: "order_number", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	_, err = r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{bson.E{Key: "customer_id", Value: 1}},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) Save(order *entity.Order) (*entity.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if order.ID != "" {
		return r.update(ctx, order)
	}

	return r.insert(ctx, order)
}

func (r *OrderRepository) insert(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	customerID, err := bson.ObjectIDFromHex(order.CustomerID)
	if err != nil {
		return nil, err
	}
	model := &models.Order{
		ID:          bson.NewObjectID(),
		OrderNumber: string(order.OrderNumber),
		CustomerID:  customerID,
		Total:       order.Total,
		Status:      string(order.Status),
		PlacedAt:    time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = r.collection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	return &entity.Order{
		ID:          model.ID.Hex(),
		OrderNumber: entity.OrderNumber(model.OrderNumber),
		CustomerID:  model.CustomerID.Hex(),
		Total:       model.Total,
		Status:      entity.OrderStatus(model.Status),
		PlacedAt:    model.PlacedAt,
		UpdatedAt:   model.UpdatedAt,
	}, nil
}

func (r *OrderRepository) update(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	objectID, err := bson.ObjectIDFromHex(order.ID)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"customer_id": order.CustomerID,
			"total":       order.Total,
			"status":      string(order.Status),
			"updated_at":  time.Now(),
		},
	}

	filter := bson.M{"_id": objectID}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return r.FindByID(order.ID)
}

func (r *OrderRepository) FindByID(id string) (*entity.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var model models.Order
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&model)
	if err != nil {
		return nil, err
	}

	return &entity.Order{
		ID:          model.ID.Hex(),
		OrderNumber: entity.OrderNumber(model.OrderNumber),
		CustomerID:  model.CustomerID.Hex(),
		Total:       model.Total,
		Status:      entity.OrderStatus(model.Status),
		PlacedAt:    model.PlacedAt,
		UpdatedAt:   model.UpdatedAt,
	}, nil
}

func (r *OrderRepository) UpdateStatus(id string, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
