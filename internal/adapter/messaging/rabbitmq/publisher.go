package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/salobato/ordermanager/internal/core/entity"
)

type RabbitMQPublisher struct {
	channel  AMQPChannel
	exchange string
}

func NewRabbitMQPublisher(ch AMQPChannel) *RabbitMQPublisher {
	return &RabbitMQPublisher{
		channel:  ch,
		exchange: "orders.exchange",
	}
}

func (p *RabbitMQPublisher) PublishOrderStatusChanged(
	ctx context.Context,
	event entity.OrderEvent,
) error {

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("Erro ao serializar evento: %w", err)
	}

	err = p.channel.PublishWithContext(
		ctx,
		p.exchange,
		"order.status.changed",
		false,
		false,
		amqp.Publishing{
			Body:        body,
			ContentType: "application/json",
		},
	)

	if err != nil {
		return fmt.Errorf("Erro ao publicar a mensagem: %w", err)
	}

	return nil
}
