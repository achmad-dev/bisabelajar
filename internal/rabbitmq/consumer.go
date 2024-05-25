package rabbitmq

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/iancoleman/strcase"
	jsoniter "github.com/json-iterator/go"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

//go:generate mockery --name IConsumer
type IConsumer[T any] interface {
	ConsumeMessage(msg interface{}, dependencies T) error
	IsConsumed(msg interface{}) bool
}

var consumedMessages []string

type Consumer[T any] struct {
	cfg          *RabbitMQConfig
	conn         *amqp.Connection
	handler      func(queue string, msg amqp.Delivery, dependencies T) error
	jaegerTracer trace.Tracer
	ctx          context.Context
}

func (c Consumer[T]) ConsumeMessage(msg interface{}, dependencies T) error {

	strName := strings.Split(runtime.FuncForPC(reflect.ValueOf(c.handler).Pointer()).Name(), ".")
	var consumerHandlerName = strName[len(strName)-1]

	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	typeName := reflect.TypeOf(msg).Name()
	snakeTypeName := strcase.ToSnake(typeName)

	err = ch.ExchangeDeclare(
		snakeTypeName, // name
		c.cfg.Kind,    // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		fmt.Sprintf("%s_%s", snakeTypeName, "queue"), // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,        // queue name
		snakeTypeName, // routing key
		snakeTypeName, // exchange
		false,
		nil)
	if err != nil {
		return err
	}

	deliveries, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	if err != nil {
		return err
	}

	go func() {

		select {
		case <-c.ctx.Done():
			defer func(ch *amqp.Channel) {
				err := ch.Close()
				if err != nil {
				}
			}(ch)
			return

		case delivery, ok := <-deliveries:
			{
				if !ok {
					return
				}

				// Extract headers

				err := c.handler(q.Name, delivery, dependencies)
				if err != nil {
				}

				consumedMessages = append(consumedMessages, snakeTypeName)

				_, span := c.jaegerTracer.Start(c.ctx, consumerHandlerName)

				h, err := jsoniter.Marshal(delivery.Headers)

				if err != nil {
				}

				span.SetAttributes(attribute.Key("message-id").String(delivery.MessageId))
				span.SetAttributes(attribute.Key("correlation-id").String(delivery.CorrelationId))
				span.SetAttributes(attribute.Key("queue").String(q.Name))
				span.SetAttributes(attribute.Key("exchange").String(delivery.Exchange))
				span.SetAttributes(attribute.Key("routing-key").String(delivery.RoutingKey))
				span.SetAttributes(attribute.Key("ack").Bool(true))
				span.SetAttributes(attribute.Key("timestamp").String(delivery.Timestamp.String()))
				span.SetAttributes(attribute.Key("body").String(string(delivery.Body)))
				span.SetAttributes(attribute.Key("headers").String(string(h)))

				// Cannot use defer inside a for loop
				time.Sleep(1 * time.Millisecond)
				span.End()

				err = delivery.Ack(false)
				if err != nil {
				}
			}
		}
	}()

	return nil
}

func (c Consumer[T]) IsConsumed(msg interface{}) bool {
	timeOutTime := 20 * time.Second
	startTime := time.Now()
	timeOutExpired := false
	isConsumed := false

	for {
		if timeOutExpired {
			return false
		}
		if isConsumed {
			return true
		}

		time.Sleep(time.Second * 2)

		typeName := reflect.TypeOf(msg).Name()
		snakeTypeName := strcase.ToSnake(typeName)

		isConsumed = linq.From(consumedMessages).Contains(snakeTypeName)

		timeOutExpired = time.Now().Sub(startTime) > timeOutTime
	}
}

func NewConsumer[T any](ctx context.Context, cfg *RabbitMQConfig, conn *amqp.Connection, jaegerTracer trace.Tracer, handler func(queue string, msg amqp.Delivery, dependencies T) error) IConsumer[T] {
	return &Consumer[T]{ctx: ctx, cfg: cfg, conn: conn, jaegerTracer: jaegerTracer, handler: handler}
}
