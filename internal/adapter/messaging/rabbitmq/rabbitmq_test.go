package rabbitmq_test

import (
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/salobato/ordermanager/internal/config"
)

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
