package interfaces

import (
	"context"

	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
)

type UseCase interface {
	HandleError(ctx context.Context, err error, request *models.Request) error
	HandlePanic(ctx context.Context, recover any, request *models.Request)
	Handle(ctx context.Context, request *models.Request) error
}
