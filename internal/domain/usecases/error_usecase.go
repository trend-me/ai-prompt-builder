package usecases

import (
	"context"
	"errors"
	"log/slog"

	"github.com/trend-me/ai-prompt-builder/internal/config/exceptions"
	"github.com/trend-me/ai-prompt-builder/internal/config/properties"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
)

func (u UseCase) HandleError(ctx context.Context, err error, request *models.Request) error {
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
		if request != nil {
			request.Error = &errParsed
			_ = u.queueOutput.Publish(ctx, request.OutputQueue, request)
		}
		return nil
	}
	return errParsed
}
