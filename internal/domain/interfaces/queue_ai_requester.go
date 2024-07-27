package interfaces

import (
	"context"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
)

type QueueAiRequester interface {
	Publish(ctx context.Context, prompt string, request *models.Request) error
}
