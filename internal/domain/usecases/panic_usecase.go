package usecases

import (
	"context"
	"fmt"
	"github.com/trend-me/ai-prompt-builder/internal/config/exceptions"
	"log/slog"
)

func (u UseCase) HandlePanic(ctx context.Context, recover any) {
	slog.ErrorContext(ctx, "useCase.HandlePanic",
		slog.String("details", "process started"),
		slog.Any("error", recover))

	errParsed := exceptions.NewUnknownError(fmt.Sprintf("%v", recover))

	_ = u.HandleError(ctx, errParsed)
}
