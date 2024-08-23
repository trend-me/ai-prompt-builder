package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/trend-me/ai-prompt-builder/internal/config/exceptions"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
)

func (u UseCase) HandlePanic(ctx context.Context, recover any, request *models.Request) {
	slog.ErrorContext(ctx, "useCase.HandlePanic",
		slog.String("details", "process started"),
		slog.Any("error", recover))

	errParsed := exceptions.NewUnknownError(fmt.Sprintf("%v", recover))

	_ = u.HandleError(ctx, errParsed, request)
}
