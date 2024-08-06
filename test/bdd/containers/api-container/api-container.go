package rabbitmq_container

import (
	"context"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/trend-me/golang-rabbitmq-lib/rabbitmq"
)

var conn rabbitmq.Connection
var queues map[string]*rabbitmq.Queue
var compose tc.ComposeStack

func Connect() error {
	c, err := tc.NewDockerCompose("docker-compose.yml")
	if err != nil {
		return err
	}
	compose = c

	return conn.Connect("rabbit", "rabbit", "rabbitmq-container", "5672")
}
func Disconnect() error {
	err := compose.Down(context.Background())
	if err != nil {
		return err
	}
	return conn.Close()
}

func CreateQueues() {
	queues = map[string]*rabbitmq.Queue{
		"ai-requester":      nil,
		"ai-prompt-builder": nil,
	}
	for k, q := range queues {
		q = rabbitmq.NewQueue(&conn, k, rabbitmq.ContentTypeJson, true, true, true)

		if err := q.Connect(); err != nil {
			fmt.Println("Error connecting to queue")
			panic(err)
		}

	}

}

func PostMessageToQueue(name string, content []byte) error {
	q := queues[name]
	if q == nil {
		return fmt.Errorf("queue %s not initialized", name)
	}
	err := q.Publish(context.Background(), content)
	if err != nil {
		return err
	}
	return nil
}

func ConsumeMessageFromQueue(name string) (err error, content []byte, headers map[string]interface{}) {
	q := queues[name]
	if q == nil {
		err = fmt.Errorf("queue %s not initialized", name)
		return
	}
	ctx := context.Background()
	err = q.Consume(ctx, func(delivery amqp091.Delivery) error {
		content = delivery.Body
		headers = delivery.Headers
		ctx.Done()
		return nil
	})
	return
}
