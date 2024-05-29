package rabbitmq

import (
	"context"
	"errors"

	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	channel *amqp091.Channel
}

func NewProducer(conn *amqp091.Connection) (*Producer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Producer{channel: ch}, nil
}

func (p *Producer) Publish(ctx context.Context, exchange_name, mqtype, routingKey string, body []byte) error {
	err := p.channel.ExchangeDeclare(
		exchange_name, // name of the exchange
		mqtype,        // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // noWait
		nil,           // arguments
	)
	if err != nil {
		return err
	}

	if p.channel.IsClosed() {
		return errors.New("channel closed")
	}
	return p.channel.Publish(
		exchange_name, // exchange
		routingKey,    // routing key
		false,         // mandatory
		false,         // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
