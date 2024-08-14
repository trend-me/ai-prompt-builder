package queue

import (
	"context"
	"encoding/json"
	"github.com/trend-me/ai-prompt-builder/internal/config/exceptions"
	"github.com/trend-me/ai-prompt-builder/internal/domain/interfaces"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
	"log/slog"
)

type (
	ConnectionAiRequester interface {
		Publish(ctx context.Context, b []byte) (err error)
		Connect() (err error)
	}

	aiRequesterMessage struct {
		PromptRoadMapConfigName        string         `json:"prompt_road_map_config_name"`
		PromptRoadMapConfigExecutionId string         `json:"prompt_road_map_config_execution_id"`
		OutputQueue                    string         `json:"output_queue"`
		Model                          string         `json:"model"`
		Prompt                         string         `json:"prompt"`
		Metadata                       map[string]any `json:"metadata"`
	}

	AiRequester struct {
		queue ConnectionAiRequester
	}
)

func (a AiRequester) Publish(ctx context.Context, prompt string, request *models.Request) error {
	slog.InfoContext(ctx, "AiRequester.Publish",
		slog.String("details", "process started"))

	b, err := json.Marshal(aiRequesterMessage{
		PromptRoadMapConfigName:        request.PromptRoadMapConfigName,
		PromptRoadMapConfigExecutionId: request.PromptRoadMapConfigExecutionId,
		OutputQueue:                    request.OutputQueue,
		Model:                          request.Model,
		Prompt:                         prompt,
		Metadata:                       request.Metadata,
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

func NewAiRequester(queue ConnectionAiRequester) interfaces.QueueAiRequester {
	return &AiRequester{
		queue: queue,
	}
}
