//go:build integration

package rabbitmq_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/salobato/ordermanager/internal/adapter/messaging/rabbitmq"
	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tcRabbit "github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

func TestRabbitMQPublisher_PublishOrderStatusChanged(t *testing.T) {
	ctx := context.Background()

	container, err := tcRabbit.Run(ctx, "rabbitmq:3.13-management-alpine")
	require.NoError(t, err)
	t.Cleanup(func() { container.Terminate(ctx) })

	amqpURL, err := container.AmqpURL(ctx)
	require.NoError(t, err)

	conn, err := amqp.Dial(amqpURL)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	ch, err := conn.Channel()
	require.NoError(t, err)
	t.Cleanup(func() { ch.Close() })

	// declare the same exchange the publisher uses
	err = ch.ExchangeDeclare(
		"orders.exchange",
		"topic", // adjust to "direct"/"fanout" if that's what you use
		true,    // durable
		false, false, false, nil,
	)
	require.NoError(t, err)

	// exclusive auto-delete queue — lives only for this test
	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	require.NoError(t, err)

	// bind the queue to the exchange with the publisher's routing key
	err = ch.QueueBind(q.Name, "order.status.changed", "orders.exchange", false, nil)
	require.NoError(t, err)

	msgs, err := ch.Consume(q.Name, "test-consumer", true, false, false, false, nil)
	require.NoError(t, err)

	publisher := rabbitmq.NewRabbitMQPublisher(ch)
	event := entity.OrderEvent{
		OrderID:     "order-1",
		OrderNumber: "ORD-2024-000001",
		Status:      entity.OrderCreated,
	}
	err = publisher.PublishOrderStatusChanged(ctx, event)
	require.NoError(t, err)

	select {
	case msg := <-msgs:
		var received entity.OrderEvent
		require.NoError(t, json.Unmarshal(msg.Body, &received))
		assert.Equal(t, event.OrderID, received.OrderID)
		assert.Equal(t, event.OrderNumber, received.OrderNumber)
		assert.Equal(t, event.Status, received.Status)
	case <-time.After(5 * time.Second):
		t.Fatal("Timed out")
	}
}
