package queues

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/trend-me/ai-prompt-builder/internal/config/exceptions"
	"github.com/trend-me/ai-prompt-builder/internal/domain/interfaces"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
)

type (
	ConnectionOutput interface {
		Publish(ctx context.Context, b []byte) (err error)
		Connect() (err error)
	}

	ConnectionOutputGetter func(queueName string) ConnectionOutput

	outputMessage struct {
		PromptRoadMapConfigName        string                `json:"prompt_road_map_config_name"`
		PromptRoadMapStep              int                   `json:"prompt_road_map_step"`
		PromptRoadMapConfigExecutionId string                `json:"prompt_road_map_config_execution_id"`
		OutputQueue                    string                `json:"output_queue"`
		Model                          string                `json:"model"`
		Error                          *exceptions.ErrorType `json:"error,omitempty"`
		Metadata                       map[string]any        `json:"metadata"`
	}

	output struct {
		connectionOutputGetter ConnectionOutputGetter
	}
)

func (o output) Publish(ctx context.Context, name string, request *models.Request) error {
	slog.InfoContext(ctx, "AiRequester.Publish",
		slog.String("details", "process started"))

	b, err := json.Marshal(outputMessage{
		PromptRoadMapConfigName:        request.PromptRoadMapConfigName,
		PromptRoadMapConfigExecutionId: request.PromptRoadMapConfigExecutionId,
		PromptRoadMapStep:              request.PromptRoadMapStep,
		OutputQueue:                    request.OutputQueue,
		Model:                          request.Model,
		Error:                          request.Error,
		Metadata:                       request.Metadata,
	})
	if err != nil {
		return exceptions.NewValidationError("error parsing ai-callback message", err.Error())
	}

	err = o.connectionOutputGetter(name).Publish(ctx, b)
	if err != nil {
		return exceptions.NewQueueError(err.Error())
	}
	return nil
}

func NewOutput(connectionOutputGetter ConnectionOutputGetter) interfaces.QueueOutput {
	return &output{connectionOutputGetter: connectionOutputGetter}
}
