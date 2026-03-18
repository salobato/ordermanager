package mongo_test

import (
	"context"
	"testing"
	"time"

	"github.com/salobato/ordermanager/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func setupTestDB(t *testing.T) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := config.Load()
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		t.Fatalf("Erro ao conectar ao Mongo: %v", err)
	}

	db := client.Database("ordertest")

	err = db.Drop(ctx)
	if err != nil {
		t.Fatalf("Erro ao limpar o banco de dados: %v", err)
	}

	return db
}
