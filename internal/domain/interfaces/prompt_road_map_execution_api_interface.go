package interfaces

import (
	"context"
)

type ApiPromptRoadMapConfigExecution interface {
	UpdateStepInExecutionById(ctx context.Context, id string, stepInExecution int) error
}
