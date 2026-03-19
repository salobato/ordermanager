package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/salobato/ordermanager/internal/adapter/api/gin/health"
	"github.com/salobato/ordermanager/internal/adapter/api/gin/order"
	"github.com/salobato/ordermanager/internal/adapter/messaging/rabbitmq"
	mongoRepo "github.com/salobato/ordermanager/internal/adapter/repository/mongo"
	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/salobato/ordermanager/internal/core/usecase"
	"github.com/salobato/ordermanager/pkg/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

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
		PlaceOrder: usecase.WithPlaceOrderLogging(
			func(input usecase.PlaceOrderInput) (*entity.Order, error) {
				return usecase.PlaceOrder(orderRepo, counterRepo, publisher, input)
			},
			logger,
		),
		UpdateOrderStatus: usecase.WithUpdateOrderStatusLogging(
			func(input usecase.UpdateOrderStatusInput) (*entity.Order, error) {
				return usecase.UpdateOrderStatus(orderRepo, publisher, input)
			},
			logger,
		),
		FindByID: usecase.WithFindByIDLogging(
			func(id string) (*entity.Order, error) {
				return usecase.FindOrderByID(orderRepo, usecase.FindOrderByIDInput{
					OrderID: id,
				})
			},
			logger,
		),
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

	router := gin.Default()

	healthHandler := health.NewHandler(health.HealthChecks{
		Mongo:    mongoCheck,
		RabbitMQ: rabbitCheck,
	})
	health.RegisterRoutes(router, healthHandler)

	orderHandler := order.NewHandler(uc)
	order.RegisterRoutes(router, orderHandler)

	logger.Info("server starting", slog.String("port", cfg.Port))
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
