package rabbitmq

import (
	"context"
	"reflect"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/iancoleman/strcase"
	jsoniter "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

//go:generate mockery --name IPublisher
type IPublisher interface {
	PublishMessage(msg interface{}) error
	IsPublished(msg interface{}) bool
}

var publishedMessages []string

type Publisher struct {
	cfg          *RabbitMQConfig
	conn         *amqp.Connection
	jaegerTracer trace.Tracer
	ctx          context.Context
}

func (p Publisher) PublishMessage(msg interface{}) error {

	data, err := jsoniter.Marshal(msg)

	if err != nil {
		return err
	}

	typeName := reflect.TypeOf(msg).Elem().Name()
	snakeTypeName := strcase.ToSnake(typeName)

	_, span := p.jaegerTracer.Start(p.ctx, typeName)
	defer span.End()

	// Inject the context in the headers

	channel, err := p.conn.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	err = channel.ExchangeDeclare(
		snakeTypeName, // name
		p.cfg.Kind,    // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		return err
	}

	publishingMsg := amqp.Publishing{
		Body:         data,
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		MessageId:    uuid.NewV4().String(),
		Timestamp:    time.Now(),
	}

	err = channel.Publish(snakeTypeName, snakeTypeName, false, false, publishingMsg)

	if err != nil {
		return err
	}

	publishedMessages = append(publishedMessages, snakeTypeName)

	span.SetAttributes(attribute.Key("message-id").String(publishingMsg.MessageId))
	span.SetAttributes(attribute.Key("correlation-id").String(publishingMsg.CorrelationId))
	span.SetAttributes(attribute.Key("exchange").String(snakeTypeName))
	span.SetAttributes(attribute.Key("kind").String(p.cfg.Kind))
	span.SetAttributes(attribute.Key("content-type").String("application/json"))
	span.SetAttributes(attribute.Key("timestamp").String(publishingMsg.Timestamp.String()))
	span.SetAttributes(attribute.Key("body").String(string(publishingMsg.Body)))

	return nil
}

func (p Publisher) IsPublished(msg interface{}) bool {

	typeName := reflect.TypeOf(msg).Name()
	snakeTypeName := strcase.ToSnake(typeName)
	isPublished := linq.From(publishedMessages).Contains(snakeTypeName)

	return isPublished
}

func NewPublisher(ctx context.Context, cfg *RabbitMQConfig, conn *amqp.Connection, jaegerTracer trace.Tracer) IPublisher {
	return &Publisher{ctx: ctx, cfg: cfg, conn: conn, jaegerTracer: jaegerTracer}
}
