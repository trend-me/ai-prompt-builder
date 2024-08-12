package interfaces

import (
	"context"

	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
)

type ApiPromptRoadMapConfigExecution interface {
	UpdatePromptRoadMapConfigExecution(ctx context.Context, promptRoadMapConfigExecution *models.PromptRoadMapConfigExecution) error
}
