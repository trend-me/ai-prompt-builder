package controllers

import (
	"context"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/trend-me/ai-prompt-builder/internal/delivery/dtos"
	"github.com/trend-me/ai-prompt-builder/internal/delivery/parsers"
	"github.com/trend-me/ai-prompt-builder/internal/delivery/validations"
	"github.com/trend-me/ai-prompt-builder/internal/domain/interfaces"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
)

type controller struct {
	useCase interfaces.UseCase
}

func (c controller) def(ctx context.Context, requestModel *models.Request) {
	if r := recover(); r != nil {
		c.useCase.HandlePanic(ctx, r, requestModel)
	}
}

func (c controller) Handle(delivery amqp.Delivery) error {
	ctx := context.Background()
	var requestModel *models.Request
	defer c.def(ctx, requestModel)

	slog.Info("controller.Handle",
		slog.String("details", "process started"),
		slog.String("body", string(delivery.Body)),
		slog.String("messageId", delivery.MessageId),
		slog.String("userId", delivery.UserId))

	var request dtos.Request
	ctx, err := parsers.ParseDeliveryJSON(ctx, &request, delivery)
	if err != nil {
		return c.useCase.HandleError(ctx, err, requestModel)
	}

	err = validations.ValidateRequest(&request)
	if err != nil {
		return c.useCase.HandleError(ctx, err, requestModel)
	}

	requestModel = &models.Request{
		PromptRoadMapConfigName:        request.PromptRoadMapConfigName,
		PromptRoadMapStep:              request.PromptRoadMapStep,
		OutputQueue:                    request.OutputQueue,
		Model:                          request.Model,
		Metadata:                       request.Metadata,
		PromptRoadMapConfigExecutionId: request.PromptRoadMapConfigExecutionId,
	}

	err = c.useCase.Handle(ctx, requestModel)
	if err != nil {
		return c.useCase.HandleError(ctx, err, requestModel)
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
