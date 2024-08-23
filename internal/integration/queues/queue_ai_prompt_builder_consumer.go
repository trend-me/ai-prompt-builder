package queues

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/trend-me/ai-prompt-builder/internal/domain/interfaces"
)

type (
	ConnectionAiPromptBuilderConsumer interface {
		Consume(ctx context.Context, handler func(delivery amqp.Delivery) error) (chan error, error)
	}

	aiPromptBuilderConsumer struct {
		queue      ConnectionAiPromptBuilderConsumer
		controller interfaces.Controller
	}
)

func (a aiPromptBuilderConsumer) Consume(ctx context.Context) (chan error, error) {
	return a.queue.Consume(ctx, a.controller.Handle)
}

func NewAiPromptBuilderConsumer(queue ConnectionAiPromptBuilderConsumer, controller interfaces.Controller) interfaces.QueueAiPromptBuilderConsumer {
	return &aiPromptBuilderConsumer{queue: queue, controller: controller}
}
