package usecases

import (
	"context"
	"errors"
	"github.com/trend-me/ai-prompt-builder/internal/config/exceptions"
	"github.com/trend-me/ai-prompt-builder/internal/config/properties"
	"log/slog"
)

func (u UseCase) HandleError(ctx context.Context, err error) error {
	slog.ErrorContext(ctx, "useCase.HandleError",
		slog.String("details", "process started"),
		slog.String("error", err.Error()))

	var errParsed exceptions.ErrorType
	if !errors.As(err, &errParsed) {
		errParsed = exceptions.NewUnknownError(err.Error())
	}

	slog.ErrorContext(ctx, "useCase.HandleError",
		slog.String("details", "processing"),
		slog.String("errorJson", string(errParsed.JSON())))

	if errParsed.Notify {
		//todo: notify
	}

	if errParsed.Abort || properties.GetCtxRetryCount(ctx) > properties.GetMaxReceiveCount() {
		return nil
	}

	return nil
}
