package interfaces

import (
	"context"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
)

type ApiPromptRoadMap interface {
	GetPromptRoadMap(ctx context.Context, promptRoadMapId string) (*models.PromptRoadMap, error)
	UpdatePromptRoadMapConfigExecution(ctx context.Context, promptRoadMapConfigExecution *models.PromptRoadMapConfigExecution) error
}
