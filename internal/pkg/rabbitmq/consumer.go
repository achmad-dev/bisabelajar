package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	channel *amqp.Channel
	log     *logrus.Entry
}

func NewConsumer(conn *amqp.Connection) (*Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Consumer{channel: ch}, nil
}

func (c *Consumer) exchangeDeclare(exchange_name, mqtype string) error {
	err := c.channel.ExchangeDeclare(
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
	return nil
}

func (c *Consumer) Consume(ctx context.Context, exchange_name, mqtype, routing_key string, onMsg func(msg amqp.Delivery) error) error {
	err := c.exchangeDeclare(exchange_name, mqtype)
	if err != nil {
		return err
	}
	q, err := c.channel.QueueDeclare(
		exchange_name, // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		amqp.Table{
			"x-queue-type":     "quorum",
			"x-delivery-limit": 5,
		}, // arguments
	)
	if err != nil {
		return err
	}

	err = c.channel.QueueBind(
		q.Name,        // queue name
		routing_key,   // routing key
		exchange_name, // exchange
		false,         // no wait
		nil,           // arg
	)
	if err != nil {
		return err
	}
	msgs, err := c.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}
	go func() {
		select {
		case <-ctx.Done():
			defer func(ch *amqp.Channel) {
				err := ch.Close()
				if err != nil {
					c.log.Errorf("error close channel %s", err.Error())
				}
			}(c.channel)
			c.log.Infof("channel closed for for queue: %s", q.Name)
			return

		case msgs, ok := <-msgs:
			{
				if !ok {
					c.log.Error()
				}
				err := msgs.Ack(false)
				if err != nil {
					c.log.Errorf("can't ack for delivery %s", string(msgs.Body))
				}
			}
		}
	}()
	return nil
}
