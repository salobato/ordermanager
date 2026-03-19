package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/salobato/ordermanager/internal/adapter/api/gin/health"
	"github.com/salobato/ordermanager/internal/adapter/api/gin/order"
	"github.com/salobato/ordermanager/internal/adapter/messaging/rabbitmq"
	mongoRepo "github.com/salobato/ordermanager/internal/adapter/repository/mongo"
	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/salobato/ordermanager/internal/core/usecase"
	"github.com/salobato/ordermanager/pkg/config"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	cfg := config.Load()

	mongoClient, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}

	db := mongoClient.Database(cfg.MongoDatabase)

	counterRepo := mongoRepo.NewCounterRepository(db)
	orderRepo := mongoRepo.NewOrderRepository(db, counterRepo)

	conn, err := amqp.Dial(cfg.RabbitMQURI)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	publisher := rabbitmq.NewRabbitMQPublisher(ch)

	uc := order.UseCases{
		PlaceOrder: func(input usecase.PlaceOrderInput) (*entity.Order, error) {
			return usecase.PlaceOrder(orderRepo, counterRepo, publisher, input)
		},
		UpdateOrderStatus: func(input usecase.UpdateOrderStatusInput) (*entity.Order, error) {
			return usecase.UpdateOrderStatus(orderRepo, publisher, input)
		},
		FindByID: func(id string) (*entity.Order, error) {
			return usecase.FindOrderByID(orderRepo, usecase.FindOrderByIDInput{
				OrderID: id,
			})
		},
	}

	mongoCheck := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return mongoClient.Ping(ctx, nil)
	}

	rabbitCheck := func() error {
		if ch == nil {
			return fmt.Errorf("channel is nil")
		}
		return ch.Qos(0, 0, false)
	}

	h := health.HealthChecks{
		Mongo:    mongoCheck,
		RabbitMQ: rabbitCheck,
	}

	router := gin.Default()

	healthHandler := health.NewHandler(h)
	health.RegisterRoutes(router, healthHandler)

	handler := order.NewHandler(uc)
	order.RegisterRoutes(router, handler)

	log.Println("Server running on port", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
