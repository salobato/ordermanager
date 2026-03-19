package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CounterRepository struct {
	collection *mongo.Collection
}

type counter struct {
	ID        string    `bson:"_id"`
	Sequence  int64     `bson:"sequence"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func NewCounterRepository(db *mongo.Database) *CounterRepository {
	return &CounterRepository{
		collection: db.Collection("counters"),
	}
}

func (r *CounterRepository) GetNextSequence(counterName string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": counterName}
	update := bson.M{
		"$inc": bson.M{"sequence": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var result counter
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return 0, fmt.Errorf("Falha para buscar próxima sequência: %w", err)
	}

	return result.Sequence, nil
}

func (r *CounterRepository) GetCurrentSequence(counterName string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result counter
	err := r.collection.FindOne(ctx, bson.M{"_id": counterName}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}

		return 0, fmt.Errorf("Falha para buscar sequência atual: %w", err)
	}
	return result.Sequence, nil
}
