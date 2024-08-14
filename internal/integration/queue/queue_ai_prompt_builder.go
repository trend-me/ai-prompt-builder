package queue

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/trend-me/ai-prompt-builder/internal/domain/interfaces"
)

type (
	ConnectionAiPromptBuilder interface {
		Consume(ctx context.Context, handler func(delivery amqp.Delivery) error) (chan error, error)
	}

	AiPromptBuilder struct {
		queue      ConnectionAiPromptBuilder
		controller interfaces.Controller
	}
)

func (a AiPromptBuilder) Consume(ctx context.Context) (chan error, error) {
	return a.queue.Consume(ctx, a.controller.Handle)
}

func NewAiPromptBuilder(queue ConnectionAiPromptBuilder, controller interfaces.Controller) interfaces.QueueAiPromptBuilder {
	return &AiPromptBuilder{queue: queue, controller: controller}
}
