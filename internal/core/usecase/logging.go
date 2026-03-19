package usecase

import (
	"log/slog"
	"time"

	"github.com/salobato/ordermanager/internal/core/entity"
)

type PlaceOrderFunc func(PlaceOrderInput) (*entity.Order, error)

func WithPlaceOrderLogging(fn PlaceOrderFunc, logger *slog.Logger) PlaceOrderFunc {
	return func(input PlaceOrderInput) (*entity.Order, error) {
		log := logger.With(
			slog.String("use_case", "PlaceOrder"),
			slog.String("customer_id", input.CustomerID),
			slog.Float64("total", input.Total),
		)

		log.Info("starting")
		start := time.Now()

		order, err := fn(input)

		duration := time.Since(start)

		if err != nil {
			log.Error("failed",
				slog.String("error", err.Error()),
				slog.Duration("duration", duration),
			)
			return nil, err
		}

		log.Info("completed",
			slog.String("order_id", order.ID),
			slog.String("order_number", order.OrderNumber.String()),
			slog.String("status", string(order.Status)),
			slog.Duration("duration", duration),
		)

		return order, nil
	}
}

type UpdateOrderStatusFunc func(UpdateOrderStatusInput) (*entity.Order, error)

func WithUpdateOrderStatusLogging(fn UpdateOrderStatusFunc, logger *slog.Logger) UpdateOrderStatusFunc {
	return func(input UpdateOrderStatusInput) (*entity.Order, error) {
		log := logger.With(
			slog.String("use_case", "UpdateOrderStatus"),
			slog.String("order_id", input.OrderID),
			slog.String("requested_status", input.Status),
		)

		log.Info("starting")
		start := time.Now()

		order, err := fn(input)

		duration := time.Since(start)

		if err != nil {
			log.Error("failed",
				slog.String("error", err.Error()),
				slog.Duration("duration", duration),
			)
			return nil, err
		}

		log.Info("completed",
			slog.String("order_number", order.OrderNumber.String()),
			slog.String("status", string(order.Status)),
			slog.Duration("duration", duration),
		)

		return order, nil
	}
}

type FindByIDFunc func(string) (*entity.Order, error)

func WithFindByIDLogging(fn FindByIDFunc, logger *slog.Logger) FindByIDFunc {
	return func(id string) (*entity.Order, error) {
		log := logger.With(
			slog.String("use_case", "FindByID"),
			slog.String("order_id", id),
		)

		log.Info("starting")
		start := time.Now()

		order, err := fn(id)

		duration := time.Since(start)

		if err != nil {
			log.Error("failed",
				slog.String("error", err.Error()),
				slog.Duration("duration", duration),
			)
			return nil, err
		}

		log.Info("completed",
			slog.String("order_number", order.OrderNumber.String()),
			slog.String("status", string(order.Status)),
			slog.Duration("duration", duration),
		)

		return order, nil
	}
}
