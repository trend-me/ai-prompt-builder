package rabbitmq_container

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/trend-me/golang-rabbitmq-lib/rabbitmq"
	"time"
)

var conn rabbitmq.Connection
var queues map[string]*rabbitmq.Queue

func Connect() error {

	err := conn.Connect("rabbit", "rabbit", "localhost", "5672")
	if err != nil {
		return err
	}

	if err = createQueues(); err != nil {
		return err
	}
	return nil
}

func Disconnect() error {
	return conn.Close()
}

func createQueues() error {
	queues = map[string]*rabbitmq.Queue{
		"ai-requester":      nil,
		"ai-prompt-builder": nil,
	}
	for k, _ := range queues {
		queues[k] = rabbitmq.NewQueue(&conn, k, rabbitmq.ContentTypeJson, true, true, true)

		if err := queues[k].Connect(); err != nil {
			fmt.Println("Error connecting to queue")
			return err
		}

	}
	return nil
}

func PostMessageToQueue(name string, content []byte) error {
	if queues[name] == nil {
		return fmt.Errorf("queue %s not initialized", name)
	}
	_ = queues[name].Connect()
	err := queues[name].Publish(context.Background(), content)
	if err != nil {
		return err
	}
	return nil
}

func ConsumeMessageFromQueue(name string) (content []byte, headers map[string]interface{}, err error) {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://rabbit:rabbit@localhost:5672/")
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	defer func(ch *amqp.Channel) {
		_ = ch.Close()
	}(ch)

	// Declare a queue (optional, depending on if it's already declared)
	q, err := ch.QueueDeclare(
		name,  // queue name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, nil, err
	}

	// Consume messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return nil, nil, err
	}

	timeout := time.After(10 * time.Second)
	select {
	case d := <-msgs:
		content = d.Body
		headers = d.Headers
	case <-timeout:
	}

	return content, headers, nil
}
