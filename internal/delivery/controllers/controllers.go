package controllers

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/trend-me/ai-prompt-builder/internal/delivery/dtos"
	"github.com/trend-me/ai-prompt-builder/internal/delivery/parser"
	"github.com/trend-me/ai-prompt-builder/internal/delivery/validations"
	"github.com/trend-me/ai-prompt-builder/internal/domain/interfaces"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
	"log/slog"
)

type controller struct {
	useCase interfaces.UseCase
}

func (c controller) def(ctx context.Context) {
	if r := recover(); r != nil {
		c.useCase.HandlePanic(ctx, r)
	}
}

func (c controller) Handle(delivery amqp.Delivery) error {
	ctx := context.Background()
	defer c.def(ctx)

	slog.Info("controller.Handle",
		slog.String("details", "process started"),
		slog.String("body", string(delivery.Body)),
		slog.String("messageId", delivery.MessageId),
		slog.String("userId", delivery.UserId))

	var request dtos.Request
	err := parser.ParseDeliveryJSON(&request, delivery)
	if err != nil {
		return c.useCase.HandleError(ctx, err)
	}

	err = validations.ValidateRequest(&request)
	if err != nil {
		return c.useCase.HandleError(ctx, err)
	}

	requestModel := models.Request{
		PromptRoadMapId: request.PromptRoadMapId,
		OutputQueue:     request.OutputQueue,
		Model:           request.Model,
		Metadata:        request.Metadata,
	}

	err = c.useCase.Handle(ctx, requestModel)
	if err != nil {
		return c.useCase.HandleError(ctx, err)
	}

	slog.Info("controller.Handle",
		slog.String("details", "process finished"))

	return nil
}

func NewController(useCase interfaces.UseCase) interfaces.Controller {
	return &controller{
		useCase: useCase,
	}
}
