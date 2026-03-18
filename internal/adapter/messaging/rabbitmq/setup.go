package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPChannel interface {
	PublishWithContext(
		ctx context.Context,
		exchange, key string,
		mandatory, immediate bool,
		msg amqp.Publishing,
	) error
}

func DeclareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"orders.exchange",
		"topic",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,
	)
}
