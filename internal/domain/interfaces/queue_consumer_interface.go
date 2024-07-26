package interfaces

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueConsumer interface {
	Consume(ctx context.Context, handler func(delivery amqp.Delivery) error) (err error)
	Connect() (err error)
}
