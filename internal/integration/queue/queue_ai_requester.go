package queue

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/trend-me/ai-prompt-builder/internal/config/exceptions"
	"github.com/trend-me/ai-prompt-builder/internal/domain/interfaces"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
	"log/slog"
)

type (
	ConnectionAiPromptBuilder interface {
		Publish(ctx context.Context, b []byte) (err error)
		Consume(ctx context.Context, handler func(delivery amqp.Delivery) error) (chan error, error)
	}

	aiRequesterMessage struct {
		PromptRoadMapConfigName  string         `json:"prompt_road_map_config_name"`
		PromptRoadMapExecutionId string         `json:"prompt_road_map_execution_id"`
		OutputQueue              string         `json:"output_queue"`
		Model                    string         `json:"model"`
		Prompt                   string         `json:"prompt"`
		Metadata                 map[string]any `json:"metadata"`
	}

	AiRequester struct {
		queue ConnectionAiPromptBuilder
	}
)

func (a AiRequester) Publish(ctx context.Context, prompt string, request *models.Request) error {
	slog.InfoContext(ctx, "AiRequester.Publish",
		slog.String("details", "process starteds"))

	b, err := json.Marshal(aiRequesterMessage{
		PromptRoadMapConfigName: request.PromptRoadMapConfigName,
		OutputQueue:             request.OutputQueue,
		Model:                   request.Model,
		Prompt:                  prompt,
		Metadata:                request.Metadata,
	})
	if err != nil {
		return exceptions.NewValidationError("error parsing ai-requester message", err.Error())
	}

	err = a.queue.Publish(ctx, b)
	if err != nil {
		return exceptions.NewQueueError(err.Error())
	}

	return nil
}

func NewAiRequester(queue ConnectionAiPromptBuilder) interfaces.QueueAiRequester {
	return &AiRequester{
		queue: queue,
	}
}
