package interfaces

import (
	"context"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
)

type UseCase interface {
	HandleError(ctx context.Context, err error) error
	HandlePanic(ctx context.Context, recover any)
	Handle(ctx context.Context, request *models.Request) error
}
