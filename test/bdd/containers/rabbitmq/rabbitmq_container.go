package rabbitmq_container

import (
	"context"
	"errors"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
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
	q := queues[name]
	if q == nil {
		err = fmt.Errorf("queue %s not initialized", name)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	done := make(chan struct{})

	go func() {
		_, err = q.Consume(ctx, func(delivery amqp091.Delivery) error {
			content = delivery.Body
			headers = delivery.Headers
			cancel()
			close(done)
			return nil
		})

		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			close(done)
		}
	}()

	select {
	case <-ctx.Done():
	case <-done:
	}

	return
}
