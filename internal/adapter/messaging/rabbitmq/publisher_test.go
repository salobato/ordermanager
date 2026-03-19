package rabbitmq_test

import (
	"context"
	"errors"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/salobato/ordermanager/internal/adapter/messaging/rabbitmq"
	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/salobato/ordermanager/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChannel struct {
	mock.Mock
}

func (m *MockChannel) PublishWithContext(
	ctx context.Context,
	exchange, key string,
	mandatory, immediate bool,
	msg amqp.Publishing,
) error {
	args := m.Called(ctx, exchange, key, mandatory, immediate, msg)
	return args.Error(0)
}

func setupRabbitMQ(t *testing.T) *amqp.Channel {
	cfg := config.Load()
	conn, err := amqp.Dial(cfg.RabbitMQURI)
	if err != nil {
		t.Fatalf("Erro ao conectar ao RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("Erro ao abrir channel: %v", err)
	}

	err = ch.ExchangeDeclare(
		"orders.exchange",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("Erro ao declarar a exchange: %v", err)
	}

	return ch
}
func TestRabbitMQPublisher_Publish(t *testing.T) {
	ch := new(MockChannel)

	publisher := rabbitmq.NewRabbitMQPublisher(ch)

	event := entity.OrderEvent{
		OrderID:     "1234",
		OrderNumber: "ORD-2026-000001",
		Status:      "criado",
	}

	ch.On("PublishWithContext",
		mock.Anything,
		"orders.exchange",
		"order.status.changed",
		false,
		false,
		mock.MatchedBy(func(msg amqp.Publishing) bool {
			return len(msg.Body) > 0
		}),
	).Return(nil)

	err := publisher.PublishOrderStatusChanged(context.Background(), event)

	assert.NoError(t, err)
	ch.AssertExpectations(t)
}

func TestRabbitMQPublisher_PublishError(t *testing.T) {
	ch := new(MockChannel)

	publisher := rabbitmq.NewRabbitMQPublisher(ch)

	event := entity.OrderEvent{
		OrderID:     "123",
		OrderNumber: "ORD-2026-000001",
		Status:      "criado",
	}

	ch.On("PublishWithContext",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(errors.New("Erro ao publicar"))

	err := publisher.PublishOrderStatusChanged(context.Background(), event)

	assert.Error(t, err)
}

func TestRabbitMQPublisher_Integration(t *testing.T) {
	ch := setupRabbitMQ(t)

	queue, err := ch.QueueDeclare(
		"test-queue",
		false,
		true,
		true,
		false,
		nil,
	)
	assert.NoError(t, err)

	err = ch.QueueBind(
		queue.Name,
		"order.status.changed",
		"orders.exchange",
		false,
		nil,
	)
	assert.NoError(t, err)

	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	assert.NoError(t, err)

	publisher := rabbitmq.NewRabbitMQPublisher(ch)

	event := entity.OrderEvent{
		OrderID:     "123",
		OrderNumber: "ORD-2026-000001",
		Status:      "criado",
	}

	err = publisher.PublishOrderStatusChanged(context.Background(), event)
	assert.NoError(t, err)

	select {
	case msg := <-msgs:
		assert.NotEmpty(t, msg.Body)

	case <-time.After(2 * time.Second):
		t.Fatal("did not receive message")
	}
}
