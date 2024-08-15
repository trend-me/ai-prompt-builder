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

func NewQueueAiPromptBuilderConsumerConnection(connection *rabbitmq.Connection) queue.ConnectionAiPromptBuilderConsumer {
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

func NewConsumer(controller interfaces.Controller, connectionAiPromptBuilderConsumer queue.ConnectionAiPromptBuilderConsumer) interfaces.QueueAiPromptBuilderConsumer {
	return queue.NewAiPromptBuilderConsumer(connectionAiPromptBuilderConsumer, controller)
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

func InitializeConsumer() (interfaces.QueueAiPromptBuilderConsumer, error) {
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
		NewQueueAiPromptBuilderConsumerConnection,
		connections.ConnectQueue,
		NewConsumer)
	return nil, nil
}
