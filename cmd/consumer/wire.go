//go:build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"github.com/trend-me/ai-prompt-builder/internal/config/connections"
	"github.com/trend-me/ai-prompt-builder/internal/config/properties"
	"github.com/trend-me/ai-prompt-builder/internal/delivery/controllers"
	"github.com/trend-me/ai-prompt-builder/internal/domain/interfaces"
	"github.com/trend-me/ai-prompt-builder/internal/domain/usecases"
	"github.com/trend-me/ai-prompt-builder/internal/integration/api"
	"github.com/trend-me/ai-prompt-builder/internal/integration/queue"
	"github.com/trend-me/golang-rabbitmq-lib/rabbitmq"
)

func NewQueueNameAiPromptBuilder(connection *rabbitmq.Connection) queue.ConnectionAiPromptBuilder {
	return rabbitmq.NewQueue(
		connection,
		properties.QueueNameAiPromptBuilder,
		rabbitmq.ContentTypeJson,
		properties.CreateQueueIfNX(),
		true,
		true,
	)
}

func NewConsumer(controller interfaces.Controller, connectionAiPromptBuilder queue.ConnectionAiPromptBuilder) func() {
	return func() {
		_ = connectionAiPromptBuilder.Consume(context.Background(), controller.Handle)
	}
}

func NewUrlApiValidation() api.UrlApiValidation {
	return properties.UrlApiValidation
}

func NewUrlApiPromptRoadMap() api.UrlApiPromptRoadMap {
	return properties.UrlApiPromptRoadMap
}

func InitializeConsumer() (func(), error) {
	wire.Build(controllers.NewController,
		usecases.NewUseCase,
		NewUrlApiPromptRoadMap,
		api.NewPromptRoadMap,
		NewUrlApiValidation,
		api.NewValidation,
		queue.NewAiRequester,
		NewQueueNameAiPromptBuilder,
		connections.ConnectQueue,
		NewConsumer)
	return nil, nil
}
