package injector

import (
	"github.com/PesquisAi/pesquisai-rabbitmq-lib/rabbitmq"
	"github.com/trend-me/ai-prompt-builder/internal/config/properties"
	"github.com/trend-me/ai-prompt-builder/internal/delivery/controllers"
	"github.com/trend-me/ai-prompt-builder/internal/domain/interfaces"
	"github.com/trend-me/ai-prompt-builder/internal/domain/usecases"
	"net/http"
)

type Dependencies struct {
	Controller                  interfaces.Controller
	QueueConnection             *rabbitmq.Connection
	UseCase                     interfaces.UseCase
	ConsumerAiOrchestratorQueue interfaces.QueueConsumer
	QueueAiOrchestrator         interfaces.Queue
}

func (d *Dependencies) Inject() *Dependencies {

	if d.Mux == nil {
		d.Mux = http.NewServeMux()
	}

	if d.QueueConnection == nil {
		d.QueueConnection = &rabbitmq.Connection{}
	}

	if d.QueueAiOrchestrator == nil {
		queue := rabbitmq.NewQueue(
			d.QueueConnection,
			properties.QueueNameAiOrchestrator,
			rabbitmq.ContentTypeJson,
			properties.CreateQueueIfNX(),
			true,
			true)
		d.ConsumerAiOrchestratorQueue = queue
		d.QueueAiOrchestrator = queue
	}

	if d.UseCase == nil {
		d.UseCase = usecases.NewUseCase(d.RequestRepository, d.ServiceFactory)
	}

	if d.Controller == nil {
		d.Controller = controllers.NewController(d.QueueGemini, d.UseCase)
	}
	return d
}

func NewDependencies() *Dependencies {
	deps := &Dependencies{}
	deps.Inject()
	return deps
}
