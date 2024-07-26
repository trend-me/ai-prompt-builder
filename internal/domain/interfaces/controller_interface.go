package interfaces

import amqp "github.com/rabbitmq/amqp091-go"

type Controller interface {
	Handle(delivery amqp.Delivery) error
}
