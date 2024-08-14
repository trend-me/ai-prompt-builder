//go:build wireinject

package injector

import (
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

func NewQueueAiPromptBuilderConnection(connection *rabbitmq.Connection) queue.ConnectionAiPromptBuilder {
	return rabbitmq.NewQueue(
		connection,
		properties.QueueNameAiPromptBuilder,
		rabbitmq.ContentTypeJson,
		properties.CreateQueueIfNX(),
		true,
		true,
	)
}

func NewQueueAiRequesterConnection(connection *rabbitmq.Connection) queue.ConnectionAiRequester {
	return rabbitmq.NewQueue(
		connection,
		properties.QueueAiRequester,
		rabbitmq.ContentTypeJson,
		properties.CreateQueueIfNX(),
		true,
		true,
	)
}

func NewConsumer(controller interfaces.Controller, connectionAiPromptBuilder queue.ConnectionAiPromptBuilder) interfaces.QueueAiPromptBuilder {
	return queue.NewAiPromptBuilder(connectionAiPromptBuilder, controller)
}

func NewUrlApiValidation() api.UrlApiValidation {
	return properties.UrlApiValidation
}

func NewUrlApiPromptRoadMapConfig() api.UrlApiPromptRoadMapConfig {
	return properties.UrlApiPromptRoadMapConfig
}

func NewUrlApiPromptRoadMapConfigExecution() api.UrlApiPromptRoadMapConfigExecution {
	return properties.UrlApiPromptRoadMapConfigExecution
}

func InitializeConsumer() (interfaces.QueueAiPromptBuilder, error) {
	wire.Build(controllers.NewController,
		usecases.NewUseCase,
		NewUrlApiPromptRoadMapConfig,
		NewUrlApiPromptRoadMapConfigExecution,
		api.NewPromptRoadMapConfig,
		api.NewPromptRoadMapConfigExecution,
		NewUrlApiValidation,
		api.NewValidation,
		queue.NewAiRequester,
		NewQueueAiRequesterConnection,
		NewQueueAiPromptBuilderConnection,
		connections.ConnectQueue,
		NewConsumer)
	return nil, nil
}
